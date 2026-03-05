package moderation

import "somegit.dev/vikingowl/mistral-go-sdk/chat"

// Request represents a text moderation request (/v1/moderations).
type Request struct {
	Model    string `json:"model"`
	Input    Input  `json:"input"`
	Metadata map[string]any `json:"metadata,omitempty"`
}

// Input can be a single string or a slice of strings.
type Input = any

// TextInput creates a single-string moderation input.
func TextInput(s string) Input { return s }

// TextsInput creates a multi-string moderation input.
func TextsInput(ss []string) Input { return ss }

// ChatRequest represents a chat moderation request (/v1/chat/moderations).
type ChatRequest struct {
	Model string               `json:"model"`
	Input ChatInput            `json:"input"`
}

// ChatInput is either []Message or [][]Message.
type ChatInput = any

// SingleChatInput creates input from a single conversation.
func SingleChatInput(msgs []chat.Message) ChatInput { return msgs }

// MultiChatInput creates input from multiple conversations.
func MultiChatInput(convs [][]chat.Message) ChatInput { return convs }

// Response is the response from moderation endpoints.
type Response struct {
	ID      string   `json:"id"`
	Model   string   `json:"model"`
	Results []Result `json:"results"`
}

// Result holds per-input moderation results.
type Result struct {
	Categories     map[string]bool    `json:"categories"`
	CategoryScores map[string]float64 `json:"category_scores"`
}
