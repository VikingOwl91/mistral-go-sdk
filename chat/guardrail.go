package chat

// ModerationLLMAction specifies the action to take when content exceeds thresholds.
type ModerationLLMAction string

const (
	ModerationActionNone  ModerationLLMAction = "none"
	ModerationActionBlock ModerationLLMAction = "block"
)

// ModerationLLMV1CategoryThresholds defines per-category score thresholds for v1 moderation.
type ModerationLLMV1CategoryThresholds struct {
	Sexual                      *float64 `json:"sexual,omitempty"`
	HateAndDiscrimination       *float64 `json:"hate_and_discrimination,omitempty"`
	ViolenceAndThreats          *float64 `json:"violence_and_threats,omitempty"`
	DangerousAndCriminalContent *float64 `json:"dangerous_and_criminal_content,omitempty"`
	Selfharm                    *float64 `json:"selfharm,omitempty"`
	Health                      *float64 `json:"health,omitempty"`
	Financial                   *float64 `json:"financial,omitempty"`
	Law                         *float64 `json:"law,omitempty"`
	PII                         *float64 `json:"pii,omitempty"`
}

// ModerationLLMV1Config configures the v1 moderation LLM guardrail.
type ModerationLLMV1Config struct {
	ModelName                string                              `json:"model_name,omitempty"`
	CustomCategoryThresholds *ModerationLLMV1CategoryThresholds  `json:"custom_category_thresholds,omitempty"`
	IgnoreOtherCategories    bool                                `json:"ignore_other_categories,omitempty"`
	Action                   ModerationLLMAction                 `json:"action,omitempty"`
}

// ModerationLLMV2CategoryThresholds defines per-category score thresholds for v2 moderation.
// V2 splits "dangerous_and_criminal_content" into separate "dangerous" and "criminal"
// categories and adds "jailbreaking".
type ModerationLLMV2CategoryThresholds struct {
	Sexual                *float64 `json:"sexual,omitempty"`
	HateAndDiscrimination *float64 `json:"hate_and_discrimination,omitempty"`
	ViolenceAndThreats    *float64 `json:"violence_and_threats,omitempty"`
	Dangerous             *float64 `json:"dangerous,omitempty"`
	Criminal              *float64 `json:"criminal,omitempty"`
	Selfharm              *float64 `json:"selfharm,omitempty"`
	Health                *float64 `json:"health,omitempty"`
	Financial             *float64 `json:"financial,omitempty"`
	Law                   *float64 `json:"law,omitempty"`
	PII                   *float64 `json:"pii,omitempty"`
	Jailbreaking          *float64 `json:"jailbreaking,omitempty"`
}

// ModerationLLMV2Config configures the v2 moderation LLM guardrail.
type ModerationLLMV2Config struct {
	ModelName                string                             `json:"model_name,omitempty"`
	CustomCategoryThresholds *ModerationLLMV2CategoryThresholds `json:"custom_category_thresholds,omitempty"`
	IgnoreOtherCategories    bool                               `json:"ignore_other_categories,omitempty"`
	Action                   ModerationLLMAction                `json:"action,omitempty"`
}

// GuardrailConfig configures moderation guardrails for requests.
type GuardrailConfig struct {
	BlockOnError    bool                   `json:"block_on_error"`
	ModerationLLMV1 *ModerationLLMV1Config `json:"moderation_llm_v1,omitempty"`
	ModerationLLMV2 *ModerationLLMV2Config `json:"moderation_llm_v2,omitempty"`
}
