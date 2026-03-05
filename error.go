package mistral

import (
	"errors"
	"fmt"
	"net/http"
)

// APIError represents an error response from the Mistral API.
type APIError struct {
	StatusCode int    `json:"-"`
	Type       string `json:"type,omitempty"`
	Message    string `json:"message"`
	Param      string `json:"param,omitempty"`
	Code       string `json:"code,omitempty"`
}

func (e *APIError) Error() string {
	if e.Type != "" {
		return fmt.Sprintf("mistral: %s: %s (status %d)", e.Type, e.Message, e.StatusCode)
	}
	return fmt.Sprintf("mistral: %s (status %d)", e.Message, e.StatusCode)
}

// IsNotFound returns true if the error is a 404 Not Found error.
func IsNotFound(err error) bool {
	var apiErr *APIError
	return errors.As(err, &apiErr) && apiErr.StatusCode == http.StatusNotFound
}

// IsRateLimit returns true if the error is a 429 Too Many Requests error.
func IsRateLimit(err error) bool {
	var apiErr *APIError
	return errors.As(err, &apiErr) && apiErr.StatusCode == http.StatusTooManyRequests
}

// IsAuth returns true if the error is a 401 Unauthorized error.
func IsAuth(err error) bool {
	var apiErr *APIError
	return errors.As(err, &apiErr) && apiErr.StatusCode == http.StatusUnauthorized
}
