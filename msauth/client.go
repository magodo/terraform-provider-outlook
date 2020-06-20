package msauth

import "context"

type Client interface {
	// ObtainToken obtains JWT token in different kinds of grant types, depends on the type implementing the interface.
	ObtainToken(ctx context.Context, authority *authority) (*Token, error)

	// RefreshToken refresh a access token via refresh token
	RefreshToken(ctx context.Context, authority *authority, refreshToken *string) (*Token, error)

	GetClientID() string
	GetClientScope() scope
}
