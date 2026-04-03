package workflow

// CodeDefinition describes a workflow's code-level interface: its input/output
// schemas, signal/query/update handlers, and execution constraints.
type CodeDefinition struct {
	InputSchema         map[string]any      `json:"input_schema"`
	OutputSchema        map[string]any      `json:"output_schema,omitempty"`
	Signals             []SignalDefinition  `json:"signals,omitempty"`
	Queries             []QueryDefinition   `json:"queries,omitempty"`
	Updates             []UpdateDefinition  `json:"updates,omitempty"`
	EnforceDeterminism  bool                `json:"enforce_determinism,omitempty"`
	ExecutionTimeout    *float64            `json:"execution_timeout,omitempty"`
}

// SignalDefinition describes a signal handler on a workflow.
type SignalDefinition struct {
	Name        string         `json:"name"`
	InputSchema map[string]any `json:"input_schema"`
	Description *string        `json:"description,omitempty"`
}

// QueryDefinition describes a query handler on a workflow.
type QueryDefinition struct {
	Name         string         `json:"name"`
	InputSchema  map[string]any `json:"input_schema"`
	Description  *string        `json:"description,omitempty"`
	OutputSchema map[string]any `json:"output_schema,omitempty"`
}

// UpdateDefinition describes an update handler on a workflow.
type UpdateDefinition struct {
	Name         string         `json:"name"`
	InputSchema  map[string]any `json:"input_schema"`
	Description  *string        `json:"description,omitempty"`
	OutputSchema map[string]any `json:"output_schema,omitempty"`
}
