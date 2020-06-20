package msauth

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type DeviceAuthorizationUserInteractionFunc func(auth DeviceAuthorization)

func defaultDeviceAuthorizationUserInteractionFunc(auth DeviceAuthorization) {
	fmt.Printf("To sign in, use a web browser to open the page %s and enter the code %s to authenticate (with in %d sec).\n",
		auth.VerificationURI, auth.UserCode, auth.ExpiresIn)
}

type clientViaDeviceFlow struct {
	client   *HTTPClient
	clientID string
	scope    scope
	f        DeviceAuthorizationUserInteractionFunc
}

func (c *clientViaDeviceFlow) ObtainToken(ctx context.Context, authority *authority) (*Token, error) {
	if c.scope.IsEmpty() {
		return nil, errors.New(`"scope" is not specified`)
	}

	// Device authorization
	body := url.Values{}
	body.Set("client_id", c.clientID)
	body.Set("scope", c.scope.String())
	req, err := c.client.NewRequestWithContext(ctx, http.MethodPost, authority.DeviceEndpoint, strings.NewReader(body.Encode()))
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

	// Polling
	body = url.Values{}
	body.Set("grant_type", "urn:ietf:params:oauth:grant-type:device_code")
	body.Set("device_code", auth.DeviceCode)
	body.Set("client_id", c.clientID)

	interval := 5
	if auth.Interval != nil {
		interval = *auth.Interval
	}
	for {
		req, err := c.client.NewRequestWithContext(ctx, http.MethodPost, authority.TokenEndpoint, strings.NewReader(body.Encode()))
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
				return c.ObtainToken(ctx, authority)
			default:
				return nil, fmt.Errorf("access token response: %s", tokenerr.String())
			}
		}
		return token, nil
	}
}

func (c *clientViaDeviceFlow) RefreshToken(ctx context.Context, authority *authority, refreshToken *string) (*Token, error) {
	if refreshToken == nil {
		return nil, errors.New("nil refresh token")
	}
	body := url.Values{}
	body.Set("grant_type", "refresh_token")
	body.Set("refresh_token", *refreshToken)
	body.Set("scope", c.scope.String())

	req, err := c.client.NewRequestWithContext(ctx, http.MethodPost, authority.TokenEndpoint, strings.NewReader(body.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/x-www-form-urlencoded")
	token, tokenerr, err := c.client.DoToken(req)
	if err != nil {
		return nil, err
	}
	if tokenerr != nil {
		return nil, errors.New(tokenerr.String())
	}
	return token, nil
}

func (c *clientViaDeviceFlow) GetClientID() string {
	return c.clientID
}

func (c *clientViaDeviceFlow) GetClientScope() scope {
	return c.scope
}

func NewClientViaDeviceFlow(client *HTTPClient, scope scope, clientID string, f DeviceAuthorizationUserInteractionFunc) Client {
	return &clientViaDeviceFlow{
		client:   client,
		clientID: clientID,
		scope:    scope,
		f:        f,
	}
}
