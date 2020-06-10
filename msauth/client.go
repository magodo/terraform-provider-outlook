package msauth

import "context"

type Client interface {
	// ObtainToken obtains JWT token in different kinds of grant types, depends on the type implementing the interface.
	ObtainToken(ctx context.Context, authority *authority) (map[string]string, error)
}
