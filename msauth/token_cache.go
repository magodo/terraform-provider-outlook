package msauth

import (
	"fmt"
	"sync"
)

type TokenCache struct {
	sync.RWMutex
	cache map[string]*Token
}

func (c *TokenCache) Insert(clientID string, authority authority, scope scope, token *Token) {
	c.Lock()
	defer c.Unlock()
	k := cacheKey(clientID, authority, scope)
	c.cache[k] = token
}

func (c *TokenCache) Get(clientID string, authority authority, scope scope) *Token {
	c.RLock()
	defer c.RUnlock()
	k := cacheKey(clientID, authority, scope)
	return c.cache[k]
}

func (c *TokenCache) Delete(clientID string, authority authority, scope scope) {
	c.RLock()
	defer c.RUnlock()
	k := cacheKey(clientID, authority, scope)
	delete(c.cache, k)
}

func cacheKey(clientID string, authority authority, scope scope) string {
	return fmt.Sprintf("%s;%s;%s", clientID, authority.AuthorizationEndpoint, scope.String())
}

func NewTokenCache() *TokenCache {
	return &TokenCache{
		cache: make(map[string]*Token, 0),
	}
}
