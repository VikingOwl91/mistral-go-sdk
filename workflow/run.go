package workflow

// Run represents a workflow run.
type Run struct {
	ID           string          `json:"id"`
	WorkflowName string          `json:"workflow_name"`
	ExecutionID  string          `json:"execution_id"`
	Status       ExecutionStatus `json:"status"`
	StartTime    string          `json:"start_time"`
	EndTime      *string         `json:"end_time,omitempty"`
}

// ListRunsResponse is the response from listing workflow runs.
type ListRunsResponse struct {
	Runs          []Run   `json:"runs"`
	NextPageToken *string `json:"next_page_token,omitempty"`
}

// RunListParams holds query parameters for listing workflow runs.
type RunListParams struct {
	WorkflowIdentifier *string
	Search             *string
	Status             *string
	PageSize           *int
	NextPageToken      *string
}
