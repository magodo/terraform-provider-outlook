package msauth

import (
	"context"
	"net/http"
	"runtime"

	"github.com/magodo/terraform-provider-outlook/version"
)

type App struct {
	client    Client
	authority *authority
}

func NewApp(client Client, authorityURL string) (*App, error) {
	c := NewHTTPClient(
		http.DefaultClient,
		map[string][]string{
			"x-client-sku": []string{"TerraformProvider.Outlook"},
			"x-client-ver": []string{version.Version},
			"x-client-os":  []string{runtime.GOOS},
			"x-client-cpu": []string{runtime.GOARCH},
		},
	)
	if authorityURL == "" {
		authorityURL = "https://login.microsoftonline.com/common"
	}
	authority, err := NewAuthority(authorityURL, c.Client)
	if err != nil {
		return nil, err
	}
	return &App{client, authority}, nil
}

func (app *App) ObtainToken(ctx context.Context) (string, error) {
	return app.client.ObtainToken(ctx, app.authority)
}
