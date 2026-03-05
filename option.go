package mistral

import (
	"net/http"
	"time"
)

// Option configures a Client.
type Option func(*Client)

// WithBaseURL sets the API base URL.
func WithBaseURL(url string) Option {
	return func(c *Client) {
		c.baseURL = url
	}
}

// WithHTTPClient sets the underlying HTTP client.
func WithHTTPClient(client *http.Client) Option {
	return func(c *Client) {
		c.httpClient = client
	}
}

// WithTimeout sets the HTTP client timeout.
func WithTimeout(d time.Duration) Option {
	return func(c *Client) {
		c.httpClient.Timeout = d
	}
}

// WithRetry configures retry behavior for transient errors.
func WithRetry(maxRetries int, baseDelay time.Duration) Option {
	return func(c *Client) {
		c.maxRetries = maxRetries
		c.retryDelay = baseDelay
	}
}
