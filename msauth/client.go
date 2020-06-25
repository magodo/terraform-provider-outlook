package msauth

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"golang.org/x/oauth2"
)

func oauthClientID(clientID, tokenURL string, scopes []string) string {
	sort.Strings(scopes)
	return fmt.Sprintf("%s @ %s (%s)", clientID, tokenURL, strings.Join(scopes, " "))
}

type Client interface {
	// ObtainTokenSource obtains token source in different kinds of grant types
	ObtainTokenSource(ctx context.Context) (oauth2.TokenSource, error)

	// ID represents a unique ID of the client from oauth's POV
	ID() string
}
