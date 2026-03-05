package chat

// CompletionResponse represents a chat completion response.
type CompletionResponse struct {
	ID      string             `json:"id"`
	Object  string             `json:"object"`
	Model   string             `json:"model"`
	Created int64              `json:"created"`
	Usage   UsageInfo          `json:"usage"`
	Choices []CompletionChoice `json:"choices"`
}

// CompletionChoice represents a single completion choice.
type CompletionChoice struct {
	Index        int              `json:"index"`
	Message      AssistantMessage `json:"message"`
	FinishReason FinishReason     `json:"finish_reason"`
}

// UsageInfo contains token usage information.
type UsageInfo struct {
	PromptTokens       int  `json:"prompt_tokens"`
	CompletionTokens   int  `json:"completion_tokens"`
	TotalTokens        int  `json:"total_tokens"`
	PromptAudioSeconds *int `json:"prompt_audio_seconds,omitempty"`
}

// CompletionChunk represents a streaming completion chunk.
type CompletionChunk struct {
	ID      string                   `json:"id"`
	Object  string                   `json:"object,omitempty"`
	Model   string                   `json:"model"`
	Created int64                    `json:"created,omitempty"`
	Usage   *UsageInfo               `json:"usage,omitempty"`
	Choices []CompletionStreamChoice `json:"choices"`
}

// CompletionStreamChoice represents a streaming completion choice.
type CompletionStreamChoice struct {
	Index        int          `json:"index"`
	Delta        DeltaMessage `json:"delta"`
	FinishReason *FinishReason `json:"finish_reason"`
}

// DeltaMessage represents a partial message in a streaming response.
type DeltaMessage struct {
	Role      string     `json:"role,omitempty"`
	Content   Content    `json:"content"`
	ToolCalls []ToolCall `json:"tool_calls,omitempty"`
}
