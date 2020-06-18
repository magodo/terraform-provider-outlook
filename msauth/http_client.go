package msauth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/hashicorp/go-retryablehttp"
)

type HTTPClient struct {
	DefaultHeader http.Header
	*retryablehttp.Client
}

func (c *HTTPClient) NewRequestWithContext(ctx context.Context, method, url string, body io.Reader) (*retryablehttp.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header = c.DefaultHeader.Clone()
	return retryablehttp.FromRequest(req)
}

// Do will send a general HTTP request and unmarshal the response into `outputPtr`.
// It returns error if the response status code is not 200.
func (c *HTTPClient) Do(req *retryablehttp.Request, outputPtr interface{}) error {
	resp, err := c.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading response body: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%s: %s ", resp.Status, string(bodyBytes))
	}

	if err := json.Unmarshal(bodyBytes, outputPtr); err != nil {
		return fmt.Errorf("unmarshalling token response: %w", err)
	}
	return nil
}

// DoToken is similar to Do, while it is specifically for access token request, in which case
// client not only care about the successful case, also needs to handle the error response.
//
// On 200, `Token` will be returned with the unmarshalled successful response.
// (as defined in: https://tools.ietf.org/html/rfc6749#section-5.1)
// On 400, `TokenError` will be returned with the unmarshalled error response.
// (as defined in: https://tools.ietf.org/html/rfc6749#section-5.2, with some possible extension,
//  e.g. https://tools.ietf.org/html/rfc8628#section-3.5)
func (c *HTTPClient) DoToken(req *retryablehttp.Request) (*Token, *TokenError, error) {
	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("reading response body: %w", err)
	}

	var (
		okbody  *Token
		errbody *TokenError
	)
	switch resp.StatusCode {
	case http.StatusOK:
		if err := json.Unmarshal(bodyBytes, &okbody); err != nil {
			return nil, nil, fmt.Errorf("unmarshalling successful token response: %w", err)
		}
	case http.StatusBadRequest:
		if err := json.Unmarshal(bodyBytes, &errbody); err != nil {
			return nil, nil, fmt.Errorf("unmarshalling error token response: %w", err)
		}
	default:
		return nil, nil, fmt.Errorf("%s: %s ", resp.Status, string(bodyBytes))
	}

	return okbody, errbody, nil
}

func NewHTTPClient(client *http.Client, defaultHeader http.Header) *HTTPClient {
	header := make(map[string][]string)
	if defaultHeader != nil {
		header = defaultHeader
	}
	retryClient := retryablehttp.NewClient()
	retryClient.HTTPClient = client
	retryClient.Logger = nil
	return &HTTPClient{
		DefaultHeader: header,
		Client:        retryClient,
	}
}
