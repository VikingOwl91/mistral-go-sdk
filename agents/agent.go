package agents

import (
	"encoding/json"
	"fmt"

	"github.com/VikingOwl91/mistral-go-sdk/chat"
)

// AgentTool is a sealed interface for agent tool types.
type AgentTool interface {
	agentToolType() string
}

// FunctionTool is a function-type tool for agents.
type FunctionTool struct {
	Type     string          `json:"type"`
	Function json.RawMessage `json:"function"`
}

func (*FunctionTool) agentToolType() string { return "function" }

// WebSearchTool enables web search capability.
type WebSearchTool struct {
	Type              string             `json:"type"`
	ToolConfiguration *ToolConfiguration `json:"tool_configuration,omitempty"`
}

func (*WebSearchTool) agentToolType() string { return "web_search" }

// WebSearchPremiumTool enables premium web search.
type WebSearchPremiumTool struct {
	Type              string             `json:"type"`
	ToolConfiguration *ToolConfiguration `json:"tool_configuration,omitempty"`
}

func (*WebSearchPremiumTool) agentToolType() string { return "web_search_premium" }

// CodeInterpreterTool enables code interpreter capability.
type CodeInterpreterTool struct {
	Type              string             `json:"type"`
	ToolConfiguration *ToolConfiguration `json:"tool_configuration,omitempty"`
}

func (*CodeInterpreterTool) agentToolType() string { return "code_interpreter" }

// ImageGenerationTool enables image generation capability.
type ImageGenerationTool struct {
	Type              string             `json:"type"`
	ToolConfiguration *ToolConfiguration `json:"tool_configuration,omitempty"`
}

func (*ImageGenerationTool) agentToolType() string { return "image_generation" }

// DocumentLibraryTool enables document library search.
type DocumentLibraryTool struct {
	Type              string             `json:"type"`
	ToolConfiguration *ToolConfiguration `json:"tool_configuration,omitempty"`
	LibraryIDs        []string           `json:"library_ids"`
}

func (*DocumentLibraryTool) agentToolType() string { return "document_library" }

// ConnectorAuth holds authorization for a custom connector.
type ConnectorAuth struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

// ConnectorTool represents a custom connector tool.
type ConnectorTool struct {
	Type              string             `json:"type"`
	ConnectorID       string             `json:"connector_id"`
	Authorization     *ConnectorAuth     `json:"authorization,omitempty"`
	ToolConfiguration *ToolConfiguration `json:"tool_configuration,omitempty"`
}

func (*ConnectorTool) agentToolType() string { return "connector" }

// UnknownAgentTool holds an unrecognized tool type.
type UnknownAgentTool struct {
	Type string
	Raw  json.RawMessage
}

func (*UnknownAgentTool) agentToolType() string { return "unknown" }

func (t *UnknownAgentTool) MarshalJSON() ([]byte, error) {
	return t.Raw, nil
}

// UnmarshalAgentTool dispatches JSON to the concrete AgentTool type
// based on the "type" discriminator field.
func UnmarshalAgentTool(data []byte) (AgentTool, error) {
	var probe struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(data, &probe); err != nil {
		return nil, fmt.Errorf("mistral: unmarshal agent tool: %w", err)
	}
	switch probe.Type {
	case "function":
		var t FunctionTool
		return &t, json.Unmarshal(data, &t)
	case "web_search":
		var t WebSearchTool
		return &t, json.Unmarshal(data, &t)
	case "web_search_premium":
		var t WebSearchPremiumTool
		return &t, json.Unmarshal(data, &t)
	case "code_interpreter":
		var t CodeInterpreterTool
		return &t, json.Unmarshal(data, &t)
	case "image_generation":
		var t ImageGenerationTool
		return &t, json.Unmarshal(data, &t)
	case "document_library":
		var t DocumentLibraryTool
		return &t, json.Unmarshal(data, &t)
	case "connector":
		var t ConnectorTool
		return &t, json.Unmarshal(data, &t)
	default:
		return &UnknownAgentTool{Type: probe.Type, Raw: json.RawMessage(data)}, nil
	}
}

// AgentTools is a slice of AgentTool with custom JSON marshaling.
type AgentTools []AgentTool

func (ts AgentTools) MarshalJSON() ([]byte, error) {
	raw := make([]json.RawMessage, len(ts))
	for i, t := range ts {
		data, err := json.Marshal(t)
		if err != nil {
			return nil, err
		}
		raw[i] = data
	}
	return json.Marshal(raw)
}

func (ts *AgentTools) UnmarshalJSON(data []byte) error {
	var raw []json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	*ts = make(AgentTools, len(raw))
	for i, r := range raw {
		tool, err := UnmarshalAgentTool(r)
		if err != nil {
			return err
		}
		(*ts)[i] = tool
	}
	return nil
}

// Agent represents a Mistral agent entity.
type Agent struct {
	ID             string            `json:"id"`
	Object         string            `json:"object"`
	Name           string            `json:"name"`
	Model          string            `json:"model"`
	Version        int               `json:"version"`
	Versions       []int             `json:"versions"`
	CreatedAt      string            `json:"created_at"`
	UpdatedAt      string            `json:"updated_at"`
	DeploymentChat bool              `json:"deployment_chat"`
	Source         string            `json:"source"`
	Instructions   *string           `json:"instructions,omitempty"`
	Description    *string           `json:"description,omitempty"`
	Tools          AgentTools        `json:"tools,omitempty"`
	CompletionArgs *CompletionArgs   `json:"completion_args,omitempty"`
	Guardrails     []chat.GuardrailConfig `json:"guardrails,omitempty"`
	Handoffs       []string          `json:"handoffs,omitempty"`
	Metadata       map[string]any    `json:"metadata,omitempty"`
	VersionMessage *string           `json:"version_message,omitempty"`
}

// CreateRequest creates a new agent.
type CreateRequest struct {
	Model          string            `json:"model"`
	Name           string            `json:"name"`
	Instructions   *string           `json:"instructions,omitempty"`
	Description    *string           `json:"description,omitempty"`
	Tools          AgentTools        `json:"tools,omitempty"`
	CompletionArgs *CompletionArgs   `json:"completion_args,omitempty"`
	Guardrails     []chat.GuardrailConfig `json:"guardrails,omitempty"`
	Handoffs       []string          `json:"handoffs,omitempty"`
	Metadata       map[string]any    `json:"metadata,omitempty"`
	VersionMessage *string           `json:"version_message,omitempty"`
}

// UpdateRequest updates an existing agent (creates a new version).
type UpdateRequest struct {
	Model          *string           `json:"model,omitempty"`
	Name           *string           `json:"name,omitempty"`
	Instructions   *string           `json:"instructions,omitempty"`
	Description    *string           `json:"description,omitempty"`
	Tools          AgentTools        `json:"tools,omitempty"`
	CompletionArgs *CompletionArgs   `json:"completion_args,omitempty"`
	Guardrails     []chat.GuardrailConfig `json:"guardrails,omitempty"`
	Handoffs       []string          `json:"handoffs,omitempty"`
	DeploymentChat *bool             `json:"deployment_chat,omitempty"`
	Metadata       map[string]any    `json:"metadata,omitempty"`
	VersionMessage *string           `json:"version_message,omitempty"`
}

// ListParams holds query parameters for listing agents.
type ListParams struct {
	Page           *int
	PageSize       *int
	DeploymentChat *bool
	Name           *string
	Search         *string
	ID             *string
}

// ListVersionsParams holds pagination for version listing.
type ListVersionsParams struct {
	Page     *int
	PageSize *int
}

// AliasResponse represents an agent version alias.
type AliasResponse struct {
	Alias     string `json:"alias"`
	Version   int    `json:"version"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// CompletionArgs are whitelisted completion API arguments for agents.
type CompletionArgs struct {
	Stop             json.RawMessage `json:"stop,omitempty"`
	PresencePenalty  *float64        `json:"presence_penalty,omitempty"`
	FrequencyPenalty *float64        `json:"frequency_penalty,omitempty"`
	Temperature      *float64        `json:"temperature,omitempty"`
	TopP             *float64        `json:"top_p,omitempty"`
	MaxTokens        *int            `json:"max_tokens,omitempty"`
	RandomSeed       *int            `json:"random_seed,omitempty"`
	ResponseFormat   json.RawMessage `json:"response_format,omitempty"`
	ToolChoice       *string         `json:"tool_choice,omitempty"`
}

// ToolConfiguration holds include/exclude/confirmation lists for tools.
type ToolConfiguration struct {
	Exclude              []string `json:"exclude,omitempty"`
	Include              []string `json:"include,omitempty"`
	RequiresConfirmation []string `json:"requires_confirmation,omitempty"`
}
