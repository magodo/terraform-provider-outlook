package msauth

import (
	"context"
	"golang.org/x/oauth2/clientcredentials"
	"golang.org/x/oauth2/microsoft"

	"golang.org/x/oauth2"
)

type ClientCredentialClient interface {
	// ObtainTokenSource obtains token source in different kinds of grant types
	ObtainTokenSource(ctx context.Context) (oauth2.TokenSource, error)
}

type clientCredentialClient struct {
	config *clientcredentials.Config
}

func (c *clientCredentialClient) ObtainTokenSource(ctx context.Context) (oauth2.TokenSource, error) {
	var err error
	ts := c.config.TokenSource(ctx)
	_, err = ts.Token()
	if err != nil {
		return nil, err
	}
	return ts, nil
}

// NOTE: The value passed for the scope parameter in this request should be the resource identifier (Application ID URI)
//       of the resource you want, affixed with the .default suffix
// 		(See https://docs.microsoft.com/en-us/graph/auth-v2-service#token-request for more details)
func NewClientCredentialClient(tenantID string, clientID, clientCredential string, scopes ...string) ClientCredentialClient {
	return &clientCredentialClient{
		config: &clientcredentials.Config{
			ClientID:     clientID,
			ClientSecret: clientCredential,
			TokenURL:     microsoft.AzureADEndpoint(tenantID).TokenURL,
			Scopes:       scopes,
			AuthStyle:    oauth2.AuthStyleInParams,
		},
	}
}
