package msauth

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/pkg/browser"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/microsoft"
)

type authorizationCodeAuth struct {
	token *oauth2.Token
	err   error
}

type clientViaAuthorizationCodeFlow struct {
	client *HTTPClient
	config *oauth2.Config
}

func (c *clientViaAuthorizationCodeFlow) ObtainTokenSource(ctx context.Context, t *oauth2.Token) (oauth2.TokenSource, error) {
	var err error
	ts := c.config.TokenSource(ctx, t)
	_, err = ts.Token()
	if err != nil {
		return nil, err
	}
	return ts, nil
}

func (c *clientViaAuthorizationCodeFlow) ID() string {
	return clientIdentifier(c.config.ClientID, c.config.Endpoint.TokenURL, c.config.Scopes)
}

func (c *clientViaAuthorizationCodeFlow) ObtainToken(ctx context.Context) (*oauth2.Token, error) {

	state := randStringRunes(5)

	// launch the http server
	ch := make(chan authorizationCodeAuth)
	closech := make(chan error)
	redirectURL, err := url.Parse(c.config.RedirectURL)
	if err != nil {
		return nil, err
	}
	srv := &http.Server{Addr: redirectURL.Host}
	go func() {
		rurl, err := url.Parse(c.config.RedirectURL)
		if err != nil {
			ch <- authorizationCodeAuth{
				err: err,
			}
			return
		}
		http.HandleFunc(rurl.Path, func(w http.ResponseWriter, r *http.Request) {
			log.Printf("[INFO] handle request")
			newstate := r.URL.Query().Get("state")
			if newstate != state {
				ch <- authorizationCodeAuth{
					err: fmt.Errorf("state diff (origin: %s; new: %s)", state, newstate),
				}
				return
			}
			code := r.URL.Query().Get("code")
			if code == "" {
				ch <- authorizationCodeAuth{
					err: errors.New("authorization code is empty"),
				}
				return
			}

			// request token
			body := url.Values{
				"grant_type":    {"authorization_code"},
				"client_id":     {c.config.ClientID},
				"client_secret": {c.config.ClientSecret},
				"scope":         {strings.Join(c.config.Scopes, " ")},
				"redirect_uri":  {c.config.RedirectURL},
				"code":          {code},
			}
			req, err := NewRequestWithContext(ctx, http.MethodPost, c.config.Endpoint.TokenURL, strings.NewReader(body.Encode()))
			if err != nil {
				ch <- authorizationCodeAuth{
					err: err,
				}
				return
			}

			token, tokenerr, err := c.client.DoToken(req)
			if tokenerr != nil {
				ch <- authorizationCodeAuth{
					err: fmt.Errorf("access token response: %s", tokenerr.String()),
				}
				return
			}

			ch <- authorizationCodeAuth{token: token.ToOauth2Token()}

			w.Write([]byte(`
<!DOCTYPE html>
<html>
<body>
<p><b>You have successfully logged in!</b></p>
You can close this window now.
</body>
</html>
`))
		})
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			ch <- authorizationCodeAuth{
				err: fmt.Errorf("Authorization Code local server failed to serve: %v", err),
			}
			return
		}
		closech <- nil
		return
	}()

	// send authorization request
	query := url.Values{
		"response_type": {"code"},
		"response_mode": {"query"},
		"client_id":     {c.config.ClientID},
		"scope":         {strings.Join(c.config.Scopes, " ")},
		"redirect_uri":  {c.config.RedirectURL},
		"state":         {state},
	}
	defer srv.Close()

	if err := browser.OpenURL(fmt.Sprintf("%s?%s", c.config.Endpoint.AuthURL, query.Encode())); err != nil {
		return nil, err
	}

	// wait for token
	select {
	case result := <-ch:
		if result.err != nil {
			return nil, result.err
		}
		// close the http server
		if err := srv.Close(); err != nil {
			return nil, fmt.Errorf("failed to close server: %w", err)
		}
		// wait for http server to quit
		select {
		case result := <-ch:
			if result.err != nil {
				return nil, fmt.Errorf("server is closing: %w", result.err)
			}
		case <-closech:
			return result.token, result.err
		}
	case <-ctx.Done():
		return nil, errors.New("time out or canceled")
	}

	return nil, errors.New("never reach here")
}

func NewClientViaAuthorizationCodeFlow(tenantID, clientID, clientSecret, redirectURL string, scopes ...string) Client {
	client := retryablehttp.NewClient()
	client.Logger = nil
	return &clientViaAuthorizationCodeFlow{
		client: NewHTTPClient(client),
		config: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			Endpoint:     microsoft.AzureADEndpoint(tenantID),
			RedirectURL:  redirectURL,
			Scopes:       scopes,
		},
	}
}
