package fim

import "encoding/json"

// CompletionRequest represents a Fill-In-the-Middle completion request.
type CompletionRequest struct {
	Model       string         `json:"model"`
	Prompt      string         `json:"prompt"`
	Suffix      *string        `json:"suffix,omitempty"`
	Temperature *float64       `json:"temperature,omitempty"`
	TopP        *float64       `json:"top_p,omitempty"`
	MaxTokens   *int           `json:"max_tokens,omitempty"`
	MinTokens   *int           `json:"min_tokens,omitempty"`
	Stop        []string       `json:"stop,omitempty"`
	RandomSeed  *int           `json:"random_seed,omitempty"`
	Metadata    map[string]any `json:"metadata,omitempty"`
	stream      bool
}

// SetStream is used internally to set the stream field.
func (r *CompletionRequest) SetStream(v bool) { r.stream = v }

func (r *CompletionRequest) MarshalJSON() ([]byte, error) {
	type Alias CompletionRequest
	return json.Marshal(&struct {
		Stream bool `json:"stream"`
		*Alias
	}{
		Stream: r.stream,
		Alias:  (*Alias)(r),
	})
}
