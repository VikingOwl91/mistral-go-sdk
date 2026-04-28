package observability

// FieldType identifies the data type of a chat-completion-event field.
type FieldType string

const (
	FieldTypeEnum      FieldType = "ENUM"
	FieldTypeText      FieldType = "TEXT"
	FieldTypeInt       FieldType = "INT"
	FieldTypeFloat     FieldType = "FLOAT"
	FieldTypeBool      FieldType = "BOOL"
	FieldTypeTimestamp FieldType = "TIMESTAMP"
	FieldTypeArray     FieldType = "ARRAY"
)

// FieldOperator is a filter operator supported on observability fields.
type FieldOperator string

const (
	FieldOperatorLT          FieldOperator = "lt"
	FieldOperatorLTE         FieldOperator = "lte"
	FieldOperatorGT          FieldOperator = "gt"
	FieldOperatorGTE         FieldOperator = "gte"
	FieldOperatorStartsWith  FieldOperator = "startswith"
	FieldOperatorIStartsWith FieldOperator = "istartswith"
	FieldOperatorEndsWith    FieldOperator = "endswith"
	FieldOperatorIEndsWith   FieldOperator = "iendswith"
	FieldOperatorContains    FieldOperator = "contains"
	FieldOperatorIContains   FieldOperator = "icontains"
	FieldOperatorMatches     FieldOperator = "matches"
	FieldOperatorNotContains FieldOperator = "notcontains"
	FieldOperatorINotContain FieldOperator = "inotcontains"
	FieldOperatorEq          FieldOperator = "eq"
	FieldOperatorNeq         FieldOperator = "neq"
	FieldOperatorIsNull      FieldOperator = "isnull"
	FieldOperatorIncludes    FieldOperator = "includes"
	FieldOperatorExcludes    FieldOperator = "excludes"
	FieldOperatorLenEq       FieldOperator = "len_eq"
)

// BaseFieldDefinition describes a searchable chat-completion-event field.
type BaseFieldDefinition struct {
	Name               string          `json:"name"`
	Label              string          `json:"label"`
	Type               FieldType       `json:"type"`
	Group              *string         `json:"group,omitempty"`
	SupportedOperators []FieldOperator `json:"supported_operators"`
}

// FieldGroup groups related field definitions for UI display.
type FieldGroup struct {
	Name  string `json:"name"`
	Label string `json:"label"`
}

// ChatCompletionFields is the response of GET /v1/observability/chat-completion-fields.
type ChatCompletionFields struct {
	FieldDefinitions []BaseFieldDefinition `json:"field_definitions"`
	FieldGroups      []FieldGroup          `json:"field_groups"`
}

// ChatCompletionFieldOptions is the response of
// GET /v1/observability/chat-completion-fields/{field_name}/options.
//
// Each option may be a string, bool, or null — preserved as raw any.
type ChatCompletionFieldOptions struct {
	Options []any `json:"options"`
}

// FieldOptionCountsRequest is the body of POST options-counts.
type FieldOptionCountsRequest struct {
	FilterParams *FilterPayload `json:"filter_params,omitempty"`
}

// FieldOptionCountItem pairs a field value with how many events have it.
type FieldOptionCountItem struct {
	Value string `json:"value"`
	Count int    `json:"count"`
}

// FieldOptionCounts is the response of POST options-counts.
type FieldOptionCounts struct {
	Counts []FieldOptionCountItem `json:"counts"`
}
