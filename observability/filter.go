package observability

import "encoding/json"

// Op is a filter comparison operator.
type Op string

const (
	OpLt          Op = "lt"
	OpLte         Op = "lte"
	OpGt          Op = "gt"
	OpGte         Op = "gte"
	OpEq          Op = "eq"
	OpNeq         Op = "neq"
	OpIsNull      Op = "isnull"
	OpStartsWith  Op = "startswith"
	OpIStartsWith Op = "istartswith"
	OpEndsWith    Op = "endswith"
	OpIEndsWith   Op = "iendswith"
	OpContains    Op = "contains"
	OpIContains   Op = "icontains"
	OpMatches     Op = "matches"
	OpNotContains Op = "notcontains"
	OpINotContains Op = "inotcontains"
	OpIncludes    Op = "includes"
	OpExcludes    Op = "excludes"
	OpLenEq       Op = "len_eq"
)

// FilterCondition is a single filter comparison.
type FilterCondition struct {
	Field string `json:"field"`
	Op    Op     `json:"op"`
	Value any    `json:"value"`
}

// FilterGroup combines filters with AND/OR logic.
// The JSON keys are uppercase "AND" / "OR".
type FilterGroup struct {
	AND []json.RawMessage `json:"AND,omitempty"`
	OR  []json.RawMessage `json:"OR,omitempty"`
}

// FilterPayload wraps the top-level filter for search operations.
// Filters can be a FilterGroup or a FilterCondition.
type FilterPayload struct {
	Filters json.RawMessage `json:"filters,omitempty"`
}

// TaskStatus is the status of an async task.
type TaskStatus string

const (
	TaskStatusRunning        TaskStatus = "RUNNING"
	TaskStatusCompleted      TaskStatus = "COMPLETED"
	TaskStatusFailed         TaskStatus = "FAILED"
	TaskStatusCanceled       TaskStatus = "CANCELED"
	TaskStatusTerminated     TaskStatus = "TERMINATED"
	TaskStatusContinuedAsNew TaskStatus = "CONTINUED_AS_NEW"
	TaskStatusTimedOut       TaskStatus = "TIMED_OUT"
	TaskStatusUnknown        TaskStatus = "UNKNOWN"
)

// PaginationParams holds common pagination query parameters.
type PaginationParams struct {
	Page     *int
	PageSize *int
}

// SearchParams holds common search query parameters.
type SearchParams struct {
	Page     *int
	PageSize *int
	Q        *string
}
