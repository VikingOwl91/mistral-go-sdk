package mistral

import (
	"net/http"
	"time"
)

// Version is the SDK version string.
const Version = "0.1.0"

const (
	defaultBaseURL = "https://api.mistral.ai"
	defaultTimeout = 120 * time.Second
)

// Client is a Mistral AI API client.
type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
	maxRetries int
	retryDelay time.Duration
}

// NewClient creates a new Mistral AI client with the given API key.
func NewClient(apiKey string, opts ...Option) *Client {
	c := &Client{
		apiKey:  apiKey,
		baseURL: defaultBaseURL,
		httpClient: &http.Client{
			Timeout: defaultTimeout,
		},
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}
