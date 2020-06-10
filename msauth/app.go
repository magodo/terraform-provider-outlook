package msauth

import "net/http"

type App struct {
	client    Client
	c         *http.Client
	authority *authority
}

func NewApp(client Client, authorityURL string) (*App, error) {
	// TODO: refer to code in: _build_client()... to add some defaults (header/body)
	c := http.DefaultClient
	if authorityURL == "" {
		authorityURL = "https://login.microsoftonline.com/common"
	}
	authority, err := NewAuthority(authorityURL, c)
	if err != nil {
		return nil, err
	}
	return &App{client, c, authority}, nil
}

func (app *App) ObtainToken() map[string]string {
	return app.client.ObtainToken(app.c, app.authority)
}
