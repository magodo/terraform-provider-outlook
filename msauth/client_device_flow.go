package msauth

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/microsoft"
)

type DeviceAuthorization struct {
	DeviceCode              string  `json:"device_code"`
	UserCode                string  `json:"user_code"`
	VerificationURI         string  `json:"verification_uri"`
	VerificationURIComplete *string `json:"verification_uri_complete"`
	ExpiresIn               int     `json:"expires_in"`
	Interval                *int    `json:"interval"`
}

type DeviceAuthorizationUserInteractionFunc func(auth DeviceAuthorization)

type clientViaDeviceFlow struct {
	client *HTTPClient
	config *oauth2.Config
	f      DeviceAuthorizationUserInteractionFunc
}

func defaultDeviceAuthorizationUserInteractionFunc(auth DeviceAuthorization) {
	fmt.Printf("To sign in, use a web browser to open the page %s and enter the code %s to authenticate (with in %d sec).\n",
		auth.VerificationURI, auth.UserCode, auth.ExpiresIn)
}

func (c *clientViaDeviceFlow) ID() string {
	return oauthClientID(c.config.ClientID, c.config.Endpoint.TokenURL, c.config.Scopes)
}

func (c *clientViaDeviceFlow) ObtainTokenSource(ctx context.Context) (oauth2.TokenSource, error) {
	// Device authorization
	body := url.Values{
		"client_id": {c.config.ClientID},
		"scope":     {strings.Join(c.config.Scopes, " ")},
	}
	deviceUrl := strings.ReplaceAll(c.config.Endpoint.TokenURL, "v2.0/token", "v2.0/devicecode")
	req, err := NewRequestWithContext(ctx, http.MethodPost, deviceUrl, strings.NewReader(body.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/x-www-form-urlencoded")
	var auth DeviceAuthorization
	if err := c.client.Do(req, &auth); err != nil {
		return nil, err
	}

	// Notify user to authorize
	f := c.f
	if f == nil {
		f = defaultDeviceAuthorizationUserInteractionFunc
	}
	f(auth)

	// Polling:
	body = url.Values{
		"grant_type":  {"urn:ietf:params:oauth:grant-type:device_code"},
		"device_code": {auth.DeviceCode},
		"client_id":   {c.config.ClientID},
	}

	interval := 5
	if auth.Interval != nil {
		interval = *auth.Interval
	}
	for {
		req, err := NewRequestWithContext(ctx, http.MethodPost, c.config.Endpoint.TokenURL, strings.NewReader(body.Encode()))
		if err != nil {
			return nil, err
		}
		token, tokenerr, err := c.client.DoToken(req)
		if err != nil {
			return nil, err
		}
		if tokenerr != nil {
			switch tokenerr.Error {
			case TokenErrorAuthorizationPending:
				time.Sleep(time.Duration(interval) * time.Second)
				continue
			case TokenErrorSlowDown:
				interval += 5
				time.Sleep(time.Duration(interval) * time.Second)
				continue
			// In case authorization request was denied by user, there is no special handling for now.
			//case TokenErrorAccessDenied:
			//	return nil, errors.New(tokenerr.String())
			case TokenErrorExpiredToken:
				return c.ObtainTokenSource(ctx)
			default:
				return nil, fmt.Errorf("access token response: %s", tokenerr.String())
			}
		}
		ts := c.config.TokenSource(ctx, token.ToOauth2Token())
		if _, err := ts.Token(); err != nil {
			return nil, err
		}
		return ts, nil
	}
}

func NewClientViaDeviceFlow(tenantID string, clientID string, f DeviceAuthorizationUserInteractionFunc, scopes ...string) Client {
	client := retryablehttp.NewClient()
	client.Logger = nil
	return &clientViaDeviceFlow{
		client: NewHTTPClient(client),
		config: &oauth2.Config{
			ClientID: clientID,
			Endpoint: microsoft.AzureADEndpoint(tenantID),
			Scopes:   scopes,
		},
		f: f,
	}
}
