package mistral

import (
	"errors"
	"fmt"
	"net/http"
	"testing"
)

func TestAPIError_Error(t *testing.T) {
	tests := []struct {
		name string
		err  APIError
		want string
	}{
		{
			name: "with type",
			err:  APIError{StatusCode: 400, Type: "invalid_request", Message: "bad param"},
			want: "mistral: invalid_request: bad param (status 400)",
		},
		{
			name: "without type",
			err:  APIError{StatusCode: 500, Message: "internal error"},
			want: "mistral: internal error (status 500)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestIsNotFound(t *testing.T) {
	apiErr := &APIError{StatusCode: http.StatusNotFound, Message: "not found"}
	if !IsNotFound(apiErr) {
		t.Error("expected true for 404")
	}
	if IsNotFound(&APIError{StatusCode: 400, Message: "bad"}) {
		t.Error("expected false for 400")
	}
	if IsNotFound(errors.New("plain error")) {
		t.Error("expected false for non-API error")
	}
}

func TestIsRateLimit(t *testing.T) {
	if !IsRateLimit(&APIError{StatusCode: http.StatusTooManyRequests, Message: "slow down"}) {
		t.Error("expected true for 429")
	}
	if IsRateLimit(&APIError{StatusCode: 200, Message: "ok"}) {
		t.Error("expected false for 200")
	}
}

func TestIsAuth(t *testing.T) {
	if !IsAuth(&APIError{StatusCode: http.StatusUnauthorized, Message: "bad key"}) {
		t.Error("expected true for 401")
	}
	if IsAuth(&APIError{StatusCode: 403, Message: "forbidden"}) {
		t.Error("expected false for 403")
	}
}

func TestAPIError_Unwrap(t *testing.T) {
	apiErr := &APIError{StatusCode: 404, Message: "not found"}
	wrapped := fmt.Errorf("context: %w", apiErr)
	if !IsNotFound(wrapped) {
		t.Error("expected IsNotFound to unwrap")
	}
	if !IsRateLimit(fmt.Errorf("wrap: %w", &APIError{StatusCode: 429, Message: "limit"})) {
		t.Error("expected IsRateLimit to unwrap")
	}
	if !IsAuth(fmt.Errorf("wrap: %w", &APIError{StatusCode: 401, Message: "auth"})) {
		t.Error("expected IsAuth to unwrap")
	}
}
