package httpclient

import (
	"io"
	"net/http"
)

// HTTPClient defines the interface for an HTTP client.
type HTTPClient interface {
	Get(url string) (*http.Response, error)
	Post(url, contentType string, body io.Reader) (*http.Response, error)
}

// defaultHTTPClient is the default implementation of HTTPClient using net/http.
type defaultHTTPClient struct {
	client *http.Client
}

var _ HTTPClient = (*defaultHTTPClient)(nil)

// NewDefaultHTTPClient creates a new HTTPClient with default settings.
func NewDefaultHTTPClient() HTTPClient {
	return &defaultHTTPClient{
		client: &http.Client{},
	}
}

// Get performs a GET request.
func (c *defaultHTTPClient) Get(url string) (*http.Response, error) {
	return c.client.Get(url)
}

// Post performs a POST request.
func (c *defaultHTTPClient) Post(url, contentType string, body io.Reader) (*http.Response, error) {
	return c.client.Post(url, contentType, body)
}
