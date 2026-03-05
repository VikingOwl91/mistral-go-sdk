package conversation

import "encoding/json"

// StartRequest starts a new conversation.
type StartRequest struct {
	Inputs           Inputs           `json:"inputs"`
	Model            string           `json:"model,omitempty"`
	AgentID          string           `json:"agent_id,omitempty"`
	AgentVersion     json.RawMessage  `json:"agent_version,omitempty"`
	Instructions     *string          `json:"instructions,omitempty"`
	Tools            []Tool           `json:"tools,omitempty"`
	CompletionArgs   *CompletionArgs  `json:"completion_args,omitempty"`
	Store            *bool            `json:"store,omitempty"`
	HandoffExecution *HandoffExecution `json:"handoff_execution,omitempty"`
	Name             *string          `json:"name,omitempty"`
	Description      *string          `json:"description,omitempty"`
	Metadata         map[string]any   `json:"metadata,omitempty"`
	stream           bool
}

func (r *StartRequest) SetStream(v bool) { r.stream = v }

func (r *StartRequest) MarshalJSON() ([]byte, error) {
	type Alias StartRequest
	return json.Marshal(&struct {
		Stream bool `json:"stream"`
		*Alias
	}{
		Stream: r.stream,
		Alias:  (*Alias)(r),
	})
}

// AppendRequest appends to an existing conversation.
type AppendRequest struct {
	Inputs            Inputs                `json:"inputs"`
	CompletionArgs    *CompletionArgs       `json:"completion_args,omitempty"`
	Store             *bool                 `json:"store,omitempty"`
	HandoffExecution  *HandoffExecution     `json:"handoff_execution,omitempty"`
	ToolConfirmations []ToolCallConfirmation `json:"tool_confirmations,omitempty"`
	stream            bool
}

func (r *AppendRequest) SetStream(v bool) { r.stream = v }

func (r *AppendRequest) MarshalJSON() ([]byte, error) {
	type Alias AppendRequest
	return json.Marshal(&struct {
		Stream bool `json:"stream"`
		*Alias
	}{
		Stream: r.stream,
		Alias:  (*Alias)(r),
	})
}

// RestartRequest restarts a conversation from a specific entry.
type RestartRequest struct {
	Inputs           Inputs           `json:"inputs"`
	FromEntryID      string           `json:"from_entry_id"`
	CompletionArgs   *CompletionArgs  `json:"completion_args,omitempty"`
	Store            *bool            `json:"store,omitempty"`
	HandoffExecution *HandoffExecution `json:"handoff_execution,omitempty"`
	AgentVersion     json.RawMessage  `json:"agent_version,omitempty"`
	Metadata         map[string]any   `json:"metadata,omitempty"`
	stream           bool
}

func (r *RestartRequest) SetStream(v bool) { r.stream = v }

func (r *RestartRequest) MarshalJSON() ([]byte, error) {
	type Alias RestartRequest
	return json.Marshal(&struct {
		Stream bool `json:"stream"`
		*Alias
	}{
		Stream: r.stream,
		Alias:  (*Alias)(r),
	})
}
