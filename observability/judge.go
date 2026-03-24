package observability

import (
	"encoding/json"
	"fmt"
)

// JudgeOutputType identifies the kind of judge output.
type JudgeOutputType string

const (
	JudgeOutputClassification JudgeOutputType = "CLASSIFICATION"
	JudgeOutputRegression     JudgeOutputType = "REGRESSION"
)

// JudgeOutput is a sealed interface for judge output configurations.
type JudgeOutput interface {
	judgeOutputType() JudgeOutputType
}

// ClassificationOutput configures a classification judge.
type ClassificationOutput struct {
	Type    JudgeOutputType            `json:"type"`
	Options []ClassificationOption     `json:"options"`
}

func (*ClassificationOutput) judgeOutputType() JudgeOutputType { return JudgeOutputClassification }

// ClassificationOption is a single option for a classification judge.
type ClassificationOption struct {
	Value       string `json:"value"`
	Description string `json:"description"`
}

// RegressionOutput configures a regression judge.
type RegressionOutput struct {
	Type           JudgeOutputType `json:"type"`
	MinDescription string          `json:"min_description"`
	MaxDescription string          `json:"max_description"`
	Min            *float64        `json:"min,omitempty"`
	Max            *float64        `json:"max,omitempty"`
}

func (*RegressionOutput) judgeOutputType() JudgeOutputType { return JudgeOutputRegression }

// UnmarshalJudgeOutput dispatches to the concrete JudgeOutput type.
func UnmarshalJudgeOutput(data []byte) (JudgeOutput, error) {
	var probe struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(data, &probe); err != nil {
		return nil, fmt.Errorf("unmarshal judge output: %w", err)
	}
	switch JudgeOutputType(probe.Type) {
	case JudgeOutputClassification:
		var o ClassificationOutput
		return &o, json.Unmarshal(data, &o)
	case JudgeOutputRegression:
		var o RegressionOutput
		return &o, json.Unmarshal(data, &o)
	default:
		return nil, fmt.Errorf("unknown judge output type: %q", probe.Type)
	}
}

// Judge represents a judge entity.
type Judge struct {
	ID           string          `json:"id"`
	CreatedAt    string          `json:"created_at"`
	UpdatedAt    string          `json:"updated_at"`
	DeletedAt    *string         `json:"deleted_at,omitempty"`
	OwnerID      string          `json:"owner_id"`
	WorkspaceID  string          `json:"workspace_id"`
	Name         string          `json:"name"`
	Description  string          `json:"description"`
	ModelName    string          `json:"model_name"`
	Output       json.RawMessage `json:"output"`
	Instructions string          `json:"instructions"`
	Tools        []string        `json:"tools,omitempty"`
}

// CreateJudgeRequest creates a new judge.
type CreateJudgeRequest struct {
	Name         string          `json:"name"`
	Description  string          `json:"description"`
	ModelName    string          `json:"model_name"`
	Output       json.RawMessage `json:"output"`
	Instructions string          `json:"instructions"`
	Tools        []string        `json:"tools"`
}

// UpdateJudgeRequest updates a judge.
type UpdateJudgeRequest struct {
	Name         string          `json:"name"`
	Description  string          `json:"description"`
	ModelName    string          `json:"model_name"`
	Output       json.RawMessage `json:"output"`
	Instructions string          `json:"instructions"`
	Tools        []string        `json:"tools"`
}

// JudgeConversationRequest is the request for live-judging a conversation.
type JudgeConversationRequest struct {
	Messages   []map[string]any `json:"messages"`
	Properties map[string]any   `json:"properties,omitempty"`
}

// ListJudgesResponse is the response from listing judges.
type ListJudgesResponse struct {
	Count    int     `json:"count"`
	Results  []Judge `json:"results,omitempty"`
	Next     *string `json:"next,omitempty"`
	Previous *string `json:"previous,omitempty"`
}
