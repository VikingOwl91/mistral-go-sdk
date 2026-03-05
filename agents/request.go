package agents

import (
	"encoding/json"

	"somegit.dev/vikingowl/mistral-go-sdk/chat"
)

// CompletionRequest represents an agents completion request.
type CompletionRequest struct {
	AgentID           string               `json:"agent_id"`
	Messages          []chat.Message       `json:"-"`
	MaxTokens         *int                 `json:"max_tokens,omitempty"`
	Stop              []string             `json:"stop,omitempty"`
	RandomSeed        *int                 `json:"random_seed,omitempty"`
	Temperature       *float64             `json:"temperature,omitempty"`
	TopP              *float64             `json:"top_p,omitempty"`
	ResponseFormat    *chat.ResponseFormat `json:"response_format,omitempty"`
	Tools             []chat.Tool          `json:"tools,omitempty"`
	ToolChoice        *chat.ToolChoice     `json:"tool_choice,omitempty"`
	PresencePenalty   *float64             `json:"presence_penalty,omitempty"`
	FrequencyPenalty  *float64             `json:"frequency_penalty,omitempty"`
	N                 *int                 `json:"n,omitempty"`
	ParallelToolCalls *bool                `json:"parallel_tool_calls,omitempty"`
	Metadata          map[string]any       `json:"metadata,omitempty"`
	Prediction        *chat.Prediction     `json:"prediction,omitempty"`
	PromptMode        *chat.PromptMode     `json:"prompt_mode,omitempty"`
	stream            bool
}

// EnableStream is used by the mistral package to enable streaming on requests.
// It is not intended for direct use by consumers.
func (r *CompletionRequest) EnableStream() { r.stream = true }

func (r *CompletionRequest) MarshalJSON() ([]byte, error) {
	type Alias CompletionRequest
	return json.Marshal(&struct {
		Messages []chat.Message `json:"messages"`
		Stream   bool           `json:"stream"`
		*Alias
	}{
		Messages: r.Messages,
		Stream:   r.stream,
		Alias:    (*Alias)(r),
	})
}
