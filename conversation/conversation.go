package conversation

import (
	"encoding/json"
	"fmt"

	"somegit.dev/vikingowl/mistral-go-sdk/chat"
)

// HandoffExecution controls tool call execution.
type HandoffExecution string

const (
	HandoffClient HandoffExecution = "client"
	HandoffServer HandoffExecution = "server"
)

// Tool represents a conversation tool.
type Tool struct {
	Type              string         `json:"type"`
	Function          *chat.Function `json:"function,omitempty"`
	LibraryIDs        []string       `json:"library_ids,omitempty"`
	ToolConfiguration *ToolConfig    `json:"tool_configuration,omitempty"`
}

// ToolConfig configures tool behavior.
type ToolConfig struct {
	Exclude              []string `json:"exclude,omitempty"`
	Include              []string `json:"include,omitempty"`
	RequiresConfirmation []string `json:"requires_confirmation,omitempty"`
}

// CompletionArgs holds optional completion parameters.
type CompletionArgs struct {
	Stop             []string             `json:"stop,omitempty"`
	PresencePenalty  *float64             `json:"presence_penalty,omitempty"`
	FrequencyPenalty *float64             `json:"frequency_penalty,omitempty"`
	Temperature      *float64             `json:"temperature,omitempty"`
	TopP             *float64             `json:"top_p,omitempty"`
	MaxTokens        *int                 `json:"max_tokens,omitempty"`
	RandomSeed       *int                 `json:"random_seed,omitempty"`
	Prediction       *chat.Prediction     `json:"prediction,omitempty"`
	ResponseFormat   *chat.ResponseFormat `json:"response_format,omitempty"`
	ToolChoice       *chat.ToolChoiceMode `json:"tool_choice,omitempty"`
}

// ToolCallConfirmation confirms or denies a pending tool call.
type ToolCallConfirmation struct {
	ToolCallID   string `json:"tool_call_id"`
	Confirmation string `json:"confirmation"` // "allow" or "deny"
}

// Inputs represents conversation inputs (text string or entry array).
type Inputs struct {
	text    *string
	entries []Entry
}

// TextInputs creates Inputs from a plain text string.
func TextInputs(s string) Inputs { return Inputs{text: &s} }

// EntryInputs creates Inputs from entry objects.
func EntryInputs(entries ...Entry) Inputs { return Inputs{entries: entries} }

func (i Inputs) MarshalJSON() ([]byte, error) {
	if i.text != nil {
		return json.Marshal(*i.text)
	}
	return json.Marshal(i.entries)
}

// UsageInfo contains conversation token usage.
type UsageInfo struct {
	PromptTokens     int            `json:"prompt_tokens"`
	CompletionTokens int            `json:"completion_tokens"`
	TotalTokens      int            `json:"total_tokens"`
	ConnectorTokens  *int           `json:"connector_tokens,omitempty"`
	Connectors       map[string]int `json:"connectors,omitempty"`
}

// Response is the response from starting, appending, or restarting a conversation.
type Response struct {
	Object         string          `json:"object"`
	ConversationID string          `json:"conversation_id"`
	Outputs        []Entry         `json:"-"`
	Usage          UsageInfo       `json:"usage"`
	Guardrails     json.RawMessage `json:"guardrails,omitempty"`
}

func (r *Response) UnmarshalJSON(data []byte) error {
	type alias Response
	var raw struct {
		*alias
		Outputs []json.RawMessage `json:"outputs"`
	}
	raw.alias = (*alias)(r)
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("mistral: unmarshal conversation response: %w", err)
	}
	r.Outputs = make([]Entry, len(raw.Outputs))
	for i, o := range raw.Outputs {
		entry, err := UnmarshalEntry(o)
		if err != nil {
			return err
		}
		r.Outputs[i] = entry
	}
	return nil
}

// Conversation represents conversation metadata.
type Conversation struct {
	Object         string          `json:"object"`
	ID             string          `json:"id"`
	CreatedAt      string          `json:"created_at"`
	UpdatedAt      string          `json:"updated_at"`
	Model          string          `json:"model,omitempty"`
	AgentID        string          `json:"agent_id,omitempty"`
	AgentVersion   json.RawMessage `json:"agent_version,omitempty"`
	Name           *string         `json:"name,omitempty"`
	Description    *string         `json:"description,omitempty"`
	Metadata       map[string]any  `json:"metadata,omitempty"`
	Instructions   *string         `json:"instructions,omitempty"`
	Tools          []Tool          `json:"tools,omitempty"`
	CompletionArgs *CompletionArgs `json:"completion_args,omitempty"`
}

// History is the response from getting conversation history.
type History struct {
	Object         string  `json:"object"`
	ConversationID string  `json:"conversation_id"`
	Entries        []Entry `json:"-"`
}

func (h *History) UnmarshalJSON(data []byte) error {
	type alias History
	var raw struct {
		*alias
		Entries []json.RawMessage `json:"entries"`
	}
	raw.alias = (*alias)(h)
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("mistral: unmarshal conversation history: %w", err)
	}
	h.Entries = make([]Entry, len(raw.Entries))
	for i, e := range raw.Entries {
		entry, err := UnmarshalEntry(e)
		if err != nil {
			return err
		}
		h.Entries[i] = entry
	}
	return nil
}

// Messages is the response from getting conversation messages.
type Messages struct {
	Object         string  `json:"object"`
	ConversationID string  `json:"conversation_id"`
	Messages       []Entry `json:"-"`
}

func (m *Messages) UnmarshalJSON(data []byte) error {
	type alias Messages
	var raw struct {
		*alias
		Messages []json.RawMessage `json:"messages"`
	}
	raw.alias = (*alias)(m)
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("mistral: unmarshal conversation messages: %w", err)
	}
	m.Messages = make([]Entry, len(raw.Messages))
	for i, msg := range raw.Messages {
		entry, err := UnmarshalEntry(msg)
		if err != nil {
			return err
		}
		m.Messages[i] = entry
	}
	return nil
}
