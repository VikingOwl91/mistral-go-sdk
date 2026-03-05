package conversation

import (
	"encoding/json"
	"fmt"
)

// Event is a sealed interface for conversation streaming events.
type Event interface {
	eventType() string
}

// ResponseStartedEvent signals the start of a conversation response.
type ResponseStartedEvent struct {
	Type           string `json:"type"`
	CreatedAt      string `json:"created_at"`
	ConversationID string `json:"conversation_id"`
}

func (*ResponseStartedEvent) eventType() string { return "conversation.response.started" }

// ResponseDoneEvent signals the completion of a conversation response.
type ResponseDoneEvent struct {
	Type      string    `json:"type"`
	CreatedAt string    `json:"created_at"`
	Usage     UsageInfo `json:"usage"`
}

func (*ResponseDoneEvent) eventType() string { return "conversation.response.done" }

// ResponseErrorEvent signals an error during conversation processing.
type ResponseErrorEvent struct {
	Type      string `json:"type"`
	CreatedAt string `json:"created_at"`
	Message   string `json:"message"`
	Code      int    `json:"code"`
}

func (*ResponseErrorEvent) eventType() string { return "conversation.response.error" }

// MessageOutputEvent contains a delta of assistant message output.
type MessageOutputEvent struct {
	Type         string          `json:"type"`
	CreatedAt    string          `json:"created_at"`
	OutputIndex  int             `json:"output_index"`
	ID           string          `json:"id"`
	ContentIndex int             `json:"content_index"`
	Content      json.RawMessage `json:"content"`
	Model        *string         `json:"model,omitempty"`
	AgentID      *string         `json:"agent_id,omitempty"`
	Role         string          `json:"role"`
}

func (*MessageOutputEvent) eventType() string { return "message.output.delta" }

// ToolExecutionStartedEvent signals the start of a tool execution.
type ToolExecutionStartedEvent struct {
	Type        string  `json:"type"`
	CreatedAt   string  `json:"created_at"`
	OutputIndex int     `json:"output_index"`
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Arguments   string  `json:"arguments"`
	Model       *string `json:"model,omitempty"`
	AgentID     *string `json:"agent_id,omitempty"`
}

func (*ToolExecutionStartedEvent) eventType() string { return "tool.execution.started" }

// ToolExecutionDeltaEvent contains a delta of tool execution output.
type ToolExecutionDeltaEvent struct {
	Type        string `json:"type"`
	CreatedAt   string `json:"created_at"`
	OutputIndex int    `json:"output_index"`
	ID          string `json:"id"`
	Name        string `json:"name"`
	Arguments   string `json:"arguments"`
}

func (*ToolExecutionDeltaEvent) eventType() string { return "tool.execution.delta" }

// ToolExecutionDoneEvent signals the completion of a tool execution.
type ToolExecutionDoneEvent struct {
	Type        string         `json:"type"`
	CreatedAt   string         `json:"created_at"`
	OutputIndex int            `json:"output_index"`
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Info        map[string]any `json:"info,omitempty"`
}

func (*ToolExecutionDoneEvent) eventType() string { return "tool.execution.done" }

// FunctionCallEvent contains a delta of a function call.
type FunctionCallEvent struct {
	Type               string  `json:"type"`
	CreatedAt          string  `json:"created_at"`
	OutputIndex        int     `json:"output_index"`
	ID                 string  `json:"id"`
	Name               string  `json:"name"`
	ToolCallID         string  `json:"tool_call_id"`
	Arguments          string  `json:"arguments"`
	ConfirmationStatus *string `json:"confirmation_status,omitempty"`
	Model              *string `json:"model,omitempty"`
	AgentID            *string `json:"agent_id,omitempty"`
}

func (*FunctionCallEvent) eventType() string { return "function.call.delta" }

// AgentHandoffStartedEvent signals the start of an agent handoff.
type AgentHandoffStartedEvent struct {
	Type              string `json:"type"`
	CreatedAt         string `json:"created_at"`
	OutputIndex       int    `json:"output_index"`
	ID                string `json:"id"`
	PreviousAgentID   string `json:"previous_agent_id"`
	PreviousAgentName string `json:"previous_agent_name"`
}

func (*AgentHandoffStartedEvent) eventType() string { return "agent.handoff.started" }

// AgentHandoffDoneEvent signals the completion of an agent handoff.
type AgentHandoffDoneEvent struct {
	Type        string `json:"type"`
	CreatedAt   string `json:"created_at"`
	OutputIndex int    `json:"output_index"`
	ID          string `json:"id"`
	NextAgentID string `json:"next_agent_id"`
	NextAgentName string `json:"next_agent_name"`
}

func (*AgentHandoffDoneEvent) eventType() string { return "agent.handoff.done" }

// UnmarshalEvent dispatches JSON to the concrete Event type
// based on the "type" discriminator field.
func UnmarshalEvent(data []byte) (Event, error) {
	var probe struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(data, &probe); err != nil {
		return nil, fmt.Errorf("mistral: unmarshal event: %w", err)
	}
	switch probe.Type {
	case "conversation.response.started":
		var e ResponseStartedEvent
		return &e, json.Unmarshal(data, &e)
	case "conversation.response.done":
		var e ResponseDoneEvent
		return &e, json.Unmarshal(data, &e)
	case "conversation.response.error":
		var e ResponseErrorEvent
		return &e, json.Unmarshal(data, &e)
	case "message.output.delta":
		var e MessageOutputEvent
		return &e, json.Unmarshal(data, &e)
	case "tool.execution.started":
		var e ToolExecutionStartedEvent
		return &e, json.Unmarshal(data, &e)
	case "tool.execution.delta":
		var e ToolExecutionDeltaEvent
		return &e, json.Unmarshal(data, &e)
	case "tool.execution.done":
		var e ToolExecutionDoneEvent
		return &e, json.Unmarshal(data, &e)
	case "function.call.delta":
		var e FunctionCallEvent
		return &e, json.Unmarshal(data, &e)
	case "agent.handoff.started":
		var e AgentHandoffStartedEvent
		return &e, json.Unmarshal(data, &e)
	case "agent.handoff.done":
		var e AgentHandoffDoneEvent
		return &e, json.Unmarshal(data, &e)
	default:
		return nil, fmt.Errorf("mistral: unknown event type: %q", probe.Type)
	}
}
