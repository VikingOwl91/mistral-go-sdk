package chat

import "encoding/json"

// PromptMode controls the reasoning behavior of the model.
type PromptMode string

const (
	PromptModeReasoning PromptMode = "reasoning"
)

// Prediction provides expected completion content for optimization.
type Prediction struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}

// CompletionRequest represents a chat completion request.
type CompletionRequest struct {
	Model             string          `json:"model"`
	Messages          []Message       `json:"-"`
	Temperature       *float64        `json:"temperature,omitempty"`
	TopP              *float64        `json:"top_p,omitempty"`
	MaxTokens         *int            `json:"max_tokens,omitempty"`
	Stop              []string        `json:"stop,omitempty"`
	RandomSeed        *int            `json:"random_seed,omitempty"`
	ResponseFormat    *ResponseFormat `json:"response_format,omitempty"`
	Tools             []Tool          `json:"tools,omitempty"`
	ToolChoice        *ToolChoice     `json:"tool_choice,omitempty"`
	PresencePenalty   *float64        `json:"presence_penalty,omitempty"`
	FrequencyPenalty  *float64        `json:"frequency_penalty,omitempty"`
	N                 *int            `json:"n,omitempty"`
	SafePrompt        bool            `json:"safe_prompt,omitempty"`
	ParallelToolCalls *bool           `json:"parallel_tool_calls,omitempty"`
	Metadata          map[string]any  `json:"metadata,omitempty"`
	Prediction        *Prediction     `json:"prediction,omitempty"`
	PromptMode        *PromptMode     `json:"prompt_mode,omitempty"`
	stream            bool
}

// SetStream is used internally to set the stream field.
func (r *CompletionRequest) SetStream(v bool) { r.stream = v }

func (r *CompletionRequest) MarshalJSON() ([]byte, error) {
	type Alias CompletionRequest
	return json.Marshal(&struct {
		Messages []Message `json:"messages"`
		Stream   bool      `json:"stream"`
		*Alias
	}{
		Messages: r.Messages,
		Stream:   r.stream,
		Alias:    (*Alias)(r),
	})
}
