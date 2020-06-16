package msauth

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
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

func (c *clientViaClientCredential) ObtainToken(ctx context.Context, authority *authority) (string, error) {

	// TODO: support client assertion: See https://tools.ietf.org/html/rfc7521#section-4.2
	if c.scope.IsEmpty() {
		return "", errors.New(`"scope" is not specified`)
	}
	body := url.Values{}
	body.Set("grant_type", "client_credentials")
	body.Set("client_id", c.clientID)
	body.Set("scope", c.scope.String())

	// Use HTTP Basic authentication scheme to authenticate client.
	// See: https://tools.ietf.org/html/rfc6749#section-2.3.1
	req, err := c.client.NewRequestWithContext(ctx, http.MethodPost, authority.TokenEndpoint, strings.NewReader(body.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Add("Accept", "application/x-www-form-urlencoded")
	req.SetBasicAuth(c.clientID, c.clientCredential)

	resp, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("reading response body: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("%s: %s ", resp.Status, string(bodyBytes))
	}

	var token AccessToken
	if err := json.Unmarshal(bodyBytes, &token); err != nil {
		return "", fmt.Errorf("unmarshalling token response: %w", err)
	}

	return token.AccessToken, nil
}

func NewClientViaClientCredential(client *HTTPClient, scope scope, clientID, clientCredential string) Client {
	return &clientViaClientCredential{
		client:           client,
		clientID:         clientID,
		clientCredential: clientCredential,
		scope:            scope,
	}
}
