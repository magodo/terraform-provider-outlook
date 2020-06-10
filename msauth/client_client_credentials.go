package msauth

import (
	"encoding/base64"
	"fmt"
	"net/http"
)

type clientViaClientCredential struct {
	clientID         string
	clientCredential string
	scope            scope
	tokenEndpoint    string
}

func (c *clientViaClientCredential) ObtainToken(client *http.Client, authority *authority) map[string]string {

	// TODO: support client assertion: See https://tools.ietf.org/html/rfc7521#section-4.2
	body := map[string]string{
		"grant_type": "client_credentials",
		"client_id":  c.clientID,
	}
	if !c.scope.IsEmpty() {
		body["scope"] = c.scope.String()
	}

	// Use HTTP Basic authentication scheme to authenticate client.
	// See: https://tools.ietf.org/html/rfc6749#section-2.3.1
	header := map[string]string{
		"Accept": "application/json",
		"Authorization": fmt.Sprintf("Basic %s",
			base64.URLEncoding.EncodeToString([]byte(
				fmt.Sprintf("%s:%s", c.clientID, c.clientCredential)))),
	}

}
