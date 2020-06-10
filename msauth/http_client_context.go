package msauth

import (
	"context"
	"io"
	"net/http"
)

type HTTPClientContext struct {
	DefaultHeader http.Header
	*http.Client
}

func (c *HTTPClientContext) NewRequestWiithContext(ctx context.Context, method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	req.Header = c.DefaultHeader.Clone()
	return req, err
}
