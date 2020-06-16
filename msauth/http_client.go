package msauth

import (
	"context"
	"io"
	"net/http"
)

type HTTPClient struct {
	DefaultHeader http.Header
	*http.Client
}

func (c *HTTPClient) NewRequestWithContext(ctx context.Context, method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	req.Header = c.DefaultHeader.Clone()
	return req, err
}

func NewHTTPClient(client *http.Client, defaultHeader http.Header) *HTTPClient {
	header := make(map[string][]string)
	if defaultHeader != nil {
		header = defaultHeader
	}
	return &HTTPClient{
		DefaultHeader: header,
		Client:        client,
	}
}
