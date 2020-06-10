package msauth

import "net/http"

type Client interface {
	// ObtainToken obtains JWT token in different kinds of grant types, depends on the type implementing the interface.
	ObtainToken(client *http.Client, authority *authority) map[string]string
}
