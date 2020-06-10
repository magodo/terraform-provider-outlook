package msauth

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"
)

type clientViaClientCredential struct {
	cctx             HTTPClientContext
	clientID         string
	clientCredential string
	scope            scope
	tokenEndpoint    string
}

func (c *clientViaClientCredential) ObtainToken(ctx context.Context, authority *authority) (map[string]string, error) {

	// TODO: support client assertion: See https://tools.ietf.org/html/rfc7521#section-4.2
	body := map[string]string{
		"grant_type": "client_credentials",
		"client_id":  c.clientID,
	}
	if !c.scope.IsEmpty() {
		body["scope"] = c.scope.String()
	}
	ebody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	// Use HTTP Basic authentication scheme to authenticate client.
	// See: https://tools.ietf.org/html/rfc6749#section-2.3.1
	req, err := c.cctx.NewRequestWiithContext(ctx, http.MethodGet, authority.TokenEndpoint, bytes.NewBuffer(ebody))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	req.SetBasicAuth(url.QueryEscape(c.clientID), url.QueryEscape(c.clientCredential))

	resp, err := c.cctx.Do(req)
	if err != nil {
		return nil, err
	}

	// TODO
}
