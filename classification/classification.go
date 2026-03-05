package classification

import "somegit.dev/vikingowl/mistral-go-sdk/chat"

// Request represents a text classification request (/v1/classifications).
type Request struct {
	Model    string `json:"model"`
	Input    Input  `json:"input"`
	Metadata map[string]any `json:"metadata,omitempty"`
}

// Input can be a single string or a slice of strings.
type Input = any

// TextInput creates a single-string classification input.
func TextInput(s string) Input { return s }

// TextsInput creates a multi-string classification input.
func TextsInput(ss []string) Input { return ss }

// ChatRequest represents a chat classification request (/v1/chat/classifications).
type ChatRequest struct {
	Model string    `json:"model"`
	Input ChatInput `json:"input"`
}

// ChatInput is a single InstructRequest or a slice of them.
type ChatInput = any

// InstructRequest represents a conversation for chat classification.
type InstructRequest struct {
	Messages []chat.Message `json:"messages"`
}

// SingleChatInput creates input from a single instruct request.
func SingleChatInput(req InstructRequest) ChatInput { return req }

// MultiChatInput creates input from multiple instruct requests.
func MultiChatInput(reqs []InstructRequest) ChatInput { return reqs }

// Response is the response from classification endpoints.
type Response struct {
	ID      string         `json:"id"`
	Model   string         `json:"model"`
	Results []TargetResult `json:"results"`
}

// TargetResult holds per-input classification results.
// The outer map key is the classification target, the inner Scores map holds category→score.
type TargetResult map[string]TargetScores

// TargetScores holds the scores for a classification target.
type TargetScores struct {
	Scores map[string]float64 `json:"scores"`
}
