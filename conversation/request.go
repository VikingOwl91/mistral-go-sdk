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

// EnableStream is used by the mistral package to enable streaming on requests.
// It is not intended for direct use by consumers.
func (r *StartRequest) EnableStream() { r.stream = true }

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

// EnableStream is used by the mistral package to enable streaming on requests.
// It is not intended for direct use by consumers.
func (r *AppendRequest) EnableStream() { r.stream = true }

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

// EnableStream is used by the mistral package to enable streaming on requests.
// It is not intended for direct use by consumers.
func (r *RestartRequest) EnableStream() { r.stream = true }

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
