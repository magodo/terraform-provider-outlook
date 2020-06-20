package msauth

import (
	"context"
	"net/http"
	"runtime"
	"time"

	"github.com/magodo/terraform-provider-outlook/version"
	"github.com/pkg/errors"
)

type App struct {
	client    Client
	authority *authority
	cache     *TokenCache
}

func NewApp(client Client, authorityURL string) (*App, error) {
	c := NewHTTPClient(
		http.DefaultClient,
		map[string][]string{
			"x-client-sku": {"TerraformProvider.Outlook"},
			"x-client-ver": {version.Version},
			"x-client-os":  {runtime.GOOS},
			"x-client-cpu": {runtime.GOARCH},
		},
	)
	if authorityURL == "" {
		authorityURL = "https://login.microsoftonline.com/common"
	}
	authority, err := NewAuthority(authorityURL, c.Client)
	if err != nil {
		return nil, err
	}
	return &App{client, authority, NewTokenCache()}, nil
}

func (app *App) ObtainToken(ctx context.Context) (*Token, error) {
	token, err := app.client.ObtainToken(ctx, app.authority)
	if err != nil {
		return nil, err
	}
	token.expiresOn = time.Now().Add(time.Duration(token.ExpiresIn) * time.Second)
	app.cache.Insert(app.client.GetClientID(), *app.authority, app.client.GetClientScope(), token)
	return token, nil
}

// ObtainTokenSilently obtains token from cache, and "refresh" it if it is expired
// (or will be soon). Hence, it must be called after `ObtainToken`
func (app *App) ObtainTokenSilently(ctx context.Context) (*Token, error) {
	token := app.cache.Get(app.client.GetClientID(), *app.authority, app.client.GetClientScope())
	if token == nil {
		return nil, errors.New("token not exists in cache, perhaps not called `ObtainToken` before?")
	}
	if token.expiresOn.Sub(time.Now()) > 9999999*time.Second {
		return token, nil
	}
	token, err := app.client.RefreshToken(ctx, app.authority, token.RefreshToken)
	if err != nil {
		return nil, err
	}
	token.expiresOn = time.Now().Add(time.Duration(token.ExpiresIn) * time.Second)
	app.cache.Insert(app.client.GetClientID(), *app.authority, app.client.GetClientScope(), token)
	return token, nil
}
