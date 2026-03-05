package conversation

import (
	"encoding/json"
	"fmt"

	"somegit.dev/vikingowl/mistral-go-sdk/chat"
)

// Entry is a sealed interface for conversation history entries.
type Entry interface {
	entryType() string
}

// MessageInputEntry represents a user or assistant input message.
type MessageInputEntry struct {
	Object      string       `json:"object"`
	ID          string       `json:"id"`
	Type        string       `json:"type"`
	CreatedAt   string       `json:"created_at"`
	CompletedAt *string      `json:"completed_at,omitempty"`
	Role        string       `json:"role"`
	Content     chat.Content `json:"content"`
	Prefix      bool         `json:"prefix,omitempty"`
}

func (*MessageInputEntry) entryType() string { return "message.input" }

// MessageOutputEntry represents an assistant output message.
type MessageOutputEntry struct {
	Object      string       `json:"object"`
	ID          string       `json:"id"`
	Type        string       `json:"type"`
	CreatedAt   string       `json:"created_at"`
	CompletedAt *string      `json:"completed_at,omitempty"`
	Role        string       `json:"role"`
	Content     chat.Content `json:"content"`
	AgentID     *string      `json:"agent_id,omitempty"`
	Model       *string      `json:"model,omitempty"`
}

func (*MessageOutputEntry) entryType() string { return "message.output" }

// FunctionCallEntry represents a function call by the model.
type FunctionCallEntry struct {
	Object             string          `json:"object"`
	ID                 string          `json:"id"`
	Type               string          `json:"type"`
	CreatedAt          string          `json:"created_at"`
	CompletedAt        *string         `json:"completed_at,omitempty"`
	ToolCallID         string          `json:"tool_call_id"`
	Name               string          `json:"name"`
	Arguments          json.RawMessage `json:"arguments"`
	ConfirmationStatus *string         `json:"confirmation_status,omitempty"`
	AgentID            *string         `json:"agent_id,omitempty"`
	Model              *string         `json:"model,omitempty"`
}

func (*FunctionCallEntry) entryType() string { return "function.call" }

// FunctionResultEntry represents a function result provided by the client.
type FunctionResultEntry struct {
	Object      string  `json:"object"`
	ID          string  `json:"id"`
	Type        string  `json:"type"`
	CreatedAt   string  `json:"created_at"`
	CompletedAt *string `json:"completed_at,omitempty"`
	ToolCallID  string  `json:"tool_call_id"`
	Result      string  `json:"result"`
}

func (*FunctionResultEntry) entryType() string { return "function.result" }

// ToolExecutionEntry represents a built-in tool execution.
type ToolExecutionEntry struct {
	Object      string         `json:"object"`
	ID          string         `json:"id"`
	Type        string         `json:"type"`
	CreatedAt   string         `json:"created_at"`
	CompletedAt *string        `json:"completed_at,omitempty"`
	Name        string         `json:"name"`
	Arguments   string         `json:"arguments"`
	Info        map[string]any `json:"info,omitempty"`
	AgentID     *string        `json:"agent_id,omitempty"`
	Model       *string        `json:"model,omitempty"`
}

func (*ToolExecutionEntry) entryType() string { return "tool.execution" }

// AgentHandoffEntry represents an agent-to-agent handoff.
type AgentHandoffEntry struct {
	Object            string  `json:"object"`
	ID                string  `json:"id"`
	Type              string  `json:"type"`
	CreatedAt         string  `json:"created_at"`
	CompletedAt       *string `json:"completed_at,omitempty"`
	PreviousAgentID   string  `json:"previous_agent_id"`
	PreviousAgentName string  `json:"previous_agent_name"`
	NextAgentID       string  `json:"next_agent_id"`
	NextAgentName     string  `json:"next_agent_name"`
}

func (*AgentHandoffEntry) entryType() string { return "agent.handoff" }

// UnknownEntry holds an entry with an unrecognized type.
// This prevents the SDK from breaking when new entry types are added.
type UnknownEntry struct {
	Type string
	Raw  json.RawMessage
}

func (*UnknownEntry) entryType() string { return "unknown" }

// UnmarshalEntry dispatches JSON to the concrete Entry type
// based on the "type" discriminator field.
func UnmarshalEntry(data []byte) (Entry, error) {
	var probe struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(data, &probe); err != nil {
		return nil, fmt.Errorf("mistral: unmarshal entry: %w", err)
	}
	switch probe.Type {
	case "message.input":
		var e MessageInputEntry
		return &e, json.Unmarshal(data, &e)
	case "message.output":
		var e MessageOutputEntry
		return &e, json.Unmarshal(data, &e)
	case "function.call":
		var e FunctionCallEntry
		return &e, json.Unmarshal(data, &e)
	case "function.result":
		var e FunctionResultEntry
		return &e, json.Unmarshal(data, &e)
	case "tool.execution":
		var e ToolExecutionEntry
		return &e, json.Unmarshal(data, &e)
	case "agent.handoff":
		var e AgentHandoffEntry
		return &e, json.Unmarshal(data, &e)
	default:
		return &UnknownEntry{Type: probe.Type, Raw: json.RawMessage(data)}, nil
	}
}
