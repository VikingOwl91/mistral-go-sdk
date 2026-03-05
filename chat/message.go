package chat

import (
	"encoding/json"
	"fmt"
)

// Message is a sealed interface for chat messages.
// Only concrete types in this package can implement it.
type Message interface {
	isMessage()
	MessageRole() string
}

// SystemMessage represents a system prompt message.
type SystemMessage struct {
	Content Content `json:"content"`
}

func (*SystemMessage) isMessage()          {}
func (*SystemMessage) MessageRole() string { return "system" }

func (m *SystemMessage) MarshalJSON() ([]byte, error) {
	type alias SystemMessage
	return json.Marshal(&struct {
		Role string `json:"role"`
		*alias
	}{
		Role:  "system",
		alias: (*alias)(m),
	})
}

// UserMessage represents a user message.
type UserMessage struct {
	Content Content `json:"content"`
}

func (*UserMessage) isMessage()          {}
func (*UserMessage) MessageRole() string { return "user" }

func (m *UserMessage) MarshalJSON() ([]byte, error) {
	type alias UserMessage
	return json.Marshal(&struct {
		Role string `json:"role"`
		*alias
	}{
		Role:  "user",
		alias: (*alias)(m),
	})
}

// AssistantMessage represents an assistant response message.
type AssistantMessage struct {
	Content   Content    `json:"content"`
	ToolCalls []ToolCall `json:"tool_calls,omitempty"`
	Prefix    bool       `json:"prefix,omitempty"`
}

func (*AssistantMessage) isMessage()          {}
func (*AssistantMessage) MessageRole() string { return "assistant" }

func (m *AssistantMessage) MarshalJSON() ([]byte, error) {
	type alias AssistantMessage
	return json.Marshal(&struct {
		Role string `json:"role"`
		*alias
	}{
		Role:  "assistant",
		alias: (*alias)(m),
	})
}

// ToolMessage represents a tool result message.
type ToolMessage struct {
	Content    Content `json:"content"`
	ToolCallID string  `json:"tool_call_id,omitempty"`
	Name       string  `json:"name,omitempty"`
}

func (*ToolMessage) isMessage()          {}
func (*ToolMessage) MessageRole() string { return "tool" }

func (m *ToolMessage) MarshalJSON() ([]byte, error) {
	type alias ToolMessage
	return json.Marshal(&struct {
		Role string `json:"role"`
		*alias
	}{
		Role:  "tool",
		alias: (*alias)(m),
	})
}

// UnmarshalMessage dispatches JSON to the concrete Message type
// based on the "role" discriminator field.
func UnmarshalMessage(data []byte) (Message, error) {
	var probe struct {
		Role string `json:"role"`
	}
	if err := json.Unmarshal(data, &probe); err != nil {
		return nil, fmt.Errorf("mistral: unmarshal message: %w", err)
	}
	switch probe.Role {
	case "system":
		var m SystemMessage
		return &m, json.Unmarshal(data, &m)
	case "user":
		var m UserMessage
		return &m, json.Unmarshal(data, &m)
	case "assistant":
		var m AssistantMessage
		return &m, json.Unmarshal(data, &m)
	case "tool":
		var m ToolMessage
		return &m, json.Unmarshal(data, &m)
	default:
		return &UnknownMessage{Role: probe.Role, Raw: json.RawMessage(data)}, nil
	}
}

// UnknownMessage holds a message with an unrecognized role.
// This prevents the SDK from breaking when new roles are added.
type UnknownMessage struct {
	Role string
	Raw  json.RawMessage
}

func (*UnknownMessage) isMessage() {}

func (m *UnknownMessage) MessageRole() string { return m.Role }

func (m *UnknownMessage) MarshalJSON() ([]byte, error) {
	return m.Raw, nil
}
