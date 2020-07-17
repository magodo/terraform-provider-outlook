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

type DeviceAuthorizationAuth struct {
	DeviceCode              string  `json:"device_code"`
	UserCode                string  `json:"user_code"`
	VerificationURI         string  `json:"verification_uri"`
	VerificationURIComplete *string `json:"verification_uri_complete"`
	ExpiresIn               int     `json:"expires_in"`
	Interval                *int    `json:"interval"`
}

type DeviceAuthorizationCallback func(auth DeviceAuthorizationAuth) error

type clientViaDeviceFlow struct {
	client *HTTPClient
	config *oauth2.Config
	f      DeviceAuthorizationCallback
}

func defaultDeviceAuthorizationCallback(auth DeviceAuthorizationAuth) error {
	fmt.Printf("To sign in, use a web browser to open the page %s and enter the code %s to authenticate (with in %d sec).\n",
		auth.VerificationURI, auth.UserCode, auth.ExpiresIn)
	return nil
}

func (c *clientViaDeviceFlow) ID() string {
	return clientIdentifier(c.config.ClientID, c.config.Endpoint.TokenURL, c.config.Scopes)
}

func (c *clientViaDeviceFlow) ObtainTokenSource(ctx context.Context, t *oauth2.Token) (oauth2.TokenSource, error) {
	ts := c.config.TokenSource(ctx, t)
	if _, err := ts.Token(); err != nil {
		return nil, err
	}
	return ts, nil
}

func (c *clientViaDeviceFlow) ObtainToken(ctx context.Context) (*oauth2.Token, error) {
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
	var auth DeviceAuthorizationAuth
	if err := c.client.Do(req, &auth); err != nil {
		return nil, err
	}

	// Notify user to authorize
	f := c.f
	if f == nil {
		f = defaultDeviceAuthorizationCallback
	}
	if err := f(auth); err != nil {
		return nil, fmt.Errorf("invoking callback: %w", err)
	}

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
				return c.ObtainToken(ctx)
			default:
				return nil, fmt.Errorf("access token response: %s", tokenerr.String())
			}
		}
		return token.ToOauth2Token(), nil
	}
}

func NewClientViaDeviceFlow(tenantID string, clientID string, f DeviceAuthorizationCallback, scopes ...string) Client {
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
