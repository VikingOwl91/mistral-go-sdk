package agents

import "encoding/json"

// Agent represents a Mistral agent entity.
type Agent struct {
	ID             string           `json:"id"`
	Object         string           `json:"object"`
	Name           string           `json:"name"`
	Model          string           `json:"model"`
	Version        int              `json:"version"`
	Versions       []int            `json:"versions"`
	CreatedAt      string           `json:"created_at"`
	UpdatedAt      string           `json:"updated_at"`
	DeploymentChat bool             `json:"deployment_chat"`
	Source         string           `json:"source"`
	Instructions   *string          `json:"instructions,omitempty"`
	Description    *string          `json:"description,omitempty"`
	Tools          []json.RawMessage `json:"tools,omitempty"`
	CompletionArgs *CompletionArgs  `json:"completion_args,omitempty"`
	Guardrails     []GuardrailConfig `json:"guardrails,omitempty"`
	Handoffs       []string         `json:"handoffs,omitempty"`
	Metadata       map[string]any   `json:"metadata,omitempty"`
	VersionMessage *string          `json:"version_message,omitempty"`
}

// CreateRequest creates a new agent.
type CreateRequest struct {
	Model          string            `json:"model"`
	Name           string            `json:"name"`
	Instructions   *string           `json:"instructions,omitempty"`
	Description    *string           `json:"description,omitempty"`
	Tools          []json.RawMessage `json:"tools,omitempty"`
	CompletionArgs *CompletionArgs   `json:"completion_args,omitempty"`
	Guardrails     []GuardrailConfig `json:"guardrails,omitempty"`
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
	Tools          []json.RawMessage `json:"tools,omitempty"`
	CompletionArgs *CompletionArgs   `json:"completion_args,omitempty"`
	Guardrails     []GuardrailConfig `json:"guardrails,omitempty"`
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

// GuardrailConfig configures moderation guardrails for an agent.
type GuardrailConfig struct {
	BlockOnError    bool                 `json:"block_on_error"`
	ModerationLLMV1 *ModerationLLMV1Config `json:"moderation_llm_v1"`
}

// ModerationLLMV1Config configures the moderation LLM guardrail.
type ModerationLLMV1Config struct {
	ModelName                string          `json:"model_name,omitempty"`
	CustomCategoryThresholds json.RawMessage `json:"custom_category_thresholds,omitempty"`
	IgnoreOtherCategories    bool            `json:"ignore_other_categories,omitempty"`
	Action                   string          `json:"action,omitempty"`
}

// ToolConfiguration holds include/exclude/confirmation lists for tools.
type ToolConfiguration struct {
	Exclude              []string `json:"exclude,omitempty"`
	Include              []string `json:"include,omitempty"`
	RequiresConfirmation []string `json:"requires_confirmation,omitempty"`
}

// FunctionTool is a function-type tool for agents.
type FunctionTool struct {
	Type     string          `json:"type"`
	Function json.RawMessage `json:"function"`
}

// WebSearchTool enables web search capability.
type WebSearchTool struct {
	Type              string             `json:"type"`
	ToolConfiguration *ToolConfiguration `json:"tool_configuration,omitempty"`
}

// WebSearchPremiumTool enables premium web search.
type WebSearchPremiumTool struct {
	Type              string             `json:"type"`
	ToolConfiguration *ToolConfiguration `json:"tool_configuration,omitempty"`
}

// CodeInterpreterTool enables code interpreter capability.
type CodeInterpreterTool struct {
	Type              string             `json:"type"`
	ToolConfiguration *ToolConfiguration `json:"tool_configuration,omitempty"`
}

// ImageGenerationTool enables image generation capability.
type ImageGenerationTool struct {
	Type              string             `json:"type"`
	ToolConfiguration *ToolConfiguration `json:"tool_configuration,omitempty"`
}

// DocumentLibraryTool enables document library search.
type DocumentLibraryTool struct {
	Type              string             `json:"type"`
	ToolConfiguration *ToolConfiguration `json:"tool_configuration,omitempty"`
	LibraryIDs        []string           `json:"library_ids"`
}
