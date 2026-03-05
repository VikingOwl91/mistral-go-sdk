package chat

import (
	"encoding/json"
	"fmt"
)

// Tool represents a tool available for the model to call.
type Tool struct {
	Type     string   `json:"type"`
	Function Function `json:"function"`
}

// Function describes a function available as a tool.
type Function struct {
	Name        string         `json:"name"`
	Description string         `json:"description,omitempty"`
	Strict      bool           `json:"strict,omitempty"`
	Parameters  map[string]any `json:"parameters"`
}

// ToolCall represents a tool call made by the model.
type ToolCall struct {
	ID       string       `json:"id"`
	Type     string       `json:"type"`
	Function FunctionCall `json:"function"`
	Index    int          `json:"index,omitempty"`
}

// FunctionCall represents a function invocation in a tool call.
type FunctionCall struct {
	Name      string `json:"name"`
	Arguments string `json:"-"`
}

func (f FunctionCall) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Name      string `json:"name"`
		Arguments string `json:"arguments"`
	}{
		Name:      f.Name,
		Arguments: f.Arguments,
	})
}

func (f *FunctionCall) UnmarshalJSON(data []byte) error {
	var raw struct {
		Name      string          `json:"name"`
		Arguments json.RawMessage `json:"arguments"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	f.Name = raw.Name
	if len(raw.Arguments) > 0 && raw.Arguments[0] == '"' {
		return json.Unmarshal(raw.Arguments, &f.Arguments)
	}
	f.Arguments = string(raw.Arguments)
	return nil
}

// FunctionName identifies a function by name (used in ToolChoice).
type FunctionName struct {
	Name string `json:"name"`
}

// ToolChoiceMode is a string enum for automatic tool choice behavior.
type ToolChoiceMode string

const (
	ToolChoiceAuto     ToolChoiceMode = "auto"
	ToolChoiceNone     ToolChoiceMode = "none"
	ToolChoiceAny      ToolChoiceMode = "any"
	ToolChoiceRequired ToolChoiceMode = "required"
)

// ToolChoice specifies which tool(s) the model should use.
// Set Mode for enum values (auto/none/any/required), or
// Function to force a specific function call.
type ToolChoice struct {
	Mode     ToolChoiceMode
	Function *FunctionName
}

func (tc ToolChoice) MarshalJSON() ([]byte, error) {
	if tc.Function != nil {
		return json.Marshal(struct {
			Type     string       `json:"type"`
			Function FunctionName `json:"function"`
		}{
			Type:     "function",
			Function: *tc.Function,
		})
	}
	return json.Marshal(string(tc.Mode))
}

func (tc *ToolChoice) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		return nil
	}
	if data[0] == '"' {
		var mode string
		if err := json.Unmarshal(data, &mode); err != nil {
			return err
		}
		tc.Mode = ToolChoiceMode(mode)
		return nil
	}
	var obj struct {
		Type     string       `json:"type"`
		Function FunctionName `json:"function"`
	}
	if err := json.Unmarshal(data, &obj); err != nil {
		return fmt.Errorf("mistral: unmarshal tool choice: %w", err)
	}
	tc.Function = &obj.Function
	return nil
}

// ResponseFormatType specifies the output format.
type ResponseFormatType string

const (
	ResponseFormatText       ResponseFormatType = "text"
	ResponseFormatJSONObject ResponseFormatType = "json_object"
	ResponseFormatJSONSchema ResponseFormatType = "json_schema"
)

// ResponseFormat controls the output format of the model.
type ResponseFormat struct {
	Type       ResponseFormatType `json:"type"`
	JsonSchema *JsonSchema        `json:"json_schema,omitempty"`
}

// JsonSchema specifies a JSON schema for structured output.
type JsonSchema struct {
	Name        string         `json:"name"`
	Description *string        `json:"description,omitempty"`
	Schema      map[string]any `json:"schema"`
	Strict      bool           `json:"strict,omitempty"`
}

// FinishReason indicates why the model stopped generating.
type FinishReason string

const (
	FinishReasonStop        FinishReason = "stop"
	FinishReasonLength      FinishReason = "length"
	FinishReasonModelLength FinishReason = "model_length"
	FinishReasonError       FinishReason = "error"
	FinishReasonToolCalls   FinishReason = "tool_calls"
)
