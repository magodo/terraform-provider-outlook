package msauth

import (
	"context"

	"golang.org/x/oauth2"
)

type PrivateClient interface {
	// ObtainTokenSource obtains token source in different kinds of grant types
	ObtainTokenSource(ctx context.Context) (oauth2.TokenSource, error)
}
