package msauth

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"sync"

	"golang.org/x/oauth2"
)

type App struct {
	tokenCache tokenCache
}

type tokenCache struct {
	cache map[string]*oauth2.Token
	mutex sync.RWMutex
}

func (c *tokenCache) Get(k string) *oauth2.Token {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.cache[k]
}

func (c *tokenCache) Set(k string, v *oauth2.Token) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.cache[k] = v
}

func (app *App) ImportCache(path string) error {
	app.tokenCache.mutex.Lock()
	defer app.tokenCache.mutex.Unlock()
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, &app.tokenCache.cache)
}

func (app *App) ExportCache(path string) error {
	app.tokenCache.mutex.Lock()
	defer app.tokenCache.mutex.Unlock()
	b, err := json.MarshalIndent(app.tokenCache.cache, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, b, 0644)
}

func (app *App) ObtainTokenSourceViaClientCredential(ctx context.Context, tenantID string, clientID, clientCredential string, scopes ...string) (oauth2.TokenSource, error) {
	return NewPrivateClientViaClientCredential(tenantID, clientID, clientCredential, scopes...).ObtainTokenSource(ctx)
}

func (app *App) ObtainTokenSourceViaDeviceFlow(ctx context.Context, tenantID string, clientID string, f DeviceAuthorizationUserInteractionFunc, scopes ...string) (oauth2.TokenSource, error) {
	client := NewPublicClientViaDeviceFlow(tenantID, clientID, f, scopes...)
	t := app.tokenCache.Get(client.ID())
	if t == nil {
		var err error
		t, err = client.ObtainToken(ctx)
		if err != nil {
			return nil, err
		}
		app.tokenCache.Set(client.ID(), t)
	}
	return client.ObtainTokenSource(ctx, t)
}

func NewApp() *App {
	return &App{
		tokenCache: tokenCache{
			cache: make(map[string]*oauth2.Token),
		},
	}
}
