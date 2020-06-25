package msauth

import (
	"context"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"golang.org/x/oauth2/microsoft"
)

type clientViaClientCredential struct {
	config *clientcredentials.Config
}

func (c *clientViaClientCredential) ObtainTokenSource(ctx context.Context) (oauth2.TokenSource, error) {
	var err error
	ts := c.config.TokenSource(ctx)
	_, err = ts.Token()
	if err != nil {
		return nil, err
	}
	return ts, nil
}

func (c *clientViaClientCredential) ID() string {
	return oauthClientID(c.config.ClientID, c.config.TokenURL, c.config.Scopes)
}

// NOTE: The value passed for the scope parameter in this request should be the resource identifier (Application ID URI)
//       of the resource you want, affixed with the .default suffix
// 		(See https://docs.microsoft.com/en-us/graph/auth-v2-service#token-request for more details)
func NewClientViaClientCredential(tenantID string, clientID, clientCredential string, scopes ...string) Client {
	return &clientViaClientCredential{
		config: &clientcredentials.Config{
			ClientID:     clientID,
			ClientSecret: clientCredential,
			TokenURL:     microsoft.AzureADEndpoint(tenantID).TokenURL,
			Scopes:       scopes,
			AuthStyle:    oauth2.AuthStyleInParams,
		},
	}
}
