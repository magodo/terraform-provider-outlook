package msauth

import (
	"context"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

type clientViaClientCredential struct {
	client           *HTTPClient
	clientID         string
	clientCredential string
	scope            scope
}

func (c *clientViaClientCredential) ObtainToken(ctx context.Context, authority *authority) (*Token, error) {

	// TODO: support client assertion: See https://tools.ietf.org/html/rfc7521#section-4.2

	if c.scope.IsEmpty() {
		return nil, errors.New(`"scope" is not specified`)
	}
	body := url.Values{}
	body.Set("grant_type", "client_credentials")
	body.Set("client_id", c.clientID)
	body.Set("scope", c.scope.String())
	req, err := c.client.NewRequestWithContext(ctx, http.MethodPost, authority.TokenEndpoint, strings.NewReader(body.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/x-www-form-urlencoded")
	// Use HTTP Basic authentication scheme to authenticate client.
	// See: https://tools.ietf.org/html/rfc6749#section-2.3.1
	req.SetBasicAuth(c.clientID, c.clientCredential)

	token, tokenerr, err := c.client.DoToken(req)
	if err != nil {
		return nil, err
	}
	if tokenerr != nil {
		return nil, errors.New(tokenerr.String())
	}
	return token, nil
}

// RefreshToken of client credential will simply invoke another access token authorization.
func (c *clientViaClientCredential) RefreshToken(ctx context.Context, authority *authority, _ *string) (*Token, error) {
	return c.ObtainToken(ctx, authority)
}

func (c *clientViaClientCredential) GetClientID() string {
	return c.clientID
}

func (c *clientViaClientCredential) GetClientScope() scope {
	return c.scope
}

func NewClientViaClientCredential(client *HTTPClient, scope scope, clientID, clientCredential string) Client {
	return &clientViaClientCredential{
		client:           client,
		clientID:         clientID,
		clientCredential: clientCredential,
		scope:            scope,
	}
}
