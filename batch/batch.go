package batch

import "encoding/json"

// JobIn is the request to create a batch job.
type JobIn struct {
	Endpoint     string            `json:"endpoint"`
	InputFiles   []string          `json:"input_files,omitempty"`
	Requests     json.RawMessage   `json:"requests,omitempty"`
	Model        *string           `json:"model,omitempty"`
	AgentID      *string           `json:"agent_id,omitempty"`
	Metadata     map[string]string `json:"metadata,omitempty"`
	TimeoutHours int               `json:"timeout_hours,omitempty"`
}

// JobOut represents a batch job.
type JobOut struct {
	ID                string            `json:"id"`
	Object            string            `json:"object,omitempty"`
	InputFiles        []string          `json:"input_files"`
	Endpoint          string            `json:"endpoint"`
	Status            string            `json:"status"`
	CreatedAt         int               `json:"created_at"`
	Model             *string           `json:"model,omitempty"`
	AgentID           *string           `json:"agent_id,omitempty"`
	OutputFile        *string           `json:"output_file,omitempty"`
	ErrorFile         *string           `json:"error_file,omitempty"`
	Errors            []Error           `json:"errors"`
	Outputs           json.RawMessage   `json:"outputs,omitempty"`
	Metadata          map[string]string `json:"metadata,omitempty"`
	TotalRequests     int               `json:"total_requests"`
	CompletedRequests int               `json:"completed_requests"`
	SucceededRequests int               `json:"succeeded_requests"`
	FailedRequests    int               `json:"failed_requests"`
	StartedAt         *int              `json:"started_at,omitempty"`
	CompletedAt       *int              `json:"completed_at,omitempty"`
}

// Error describes an error encountered during batch processing.
type Error struct {
	Message string `json:"message"`
	Count   int    `json:"count"`
}

// JobsOut is a paginated list of batch jobs.
type JobsOut struct {
	Data   []JobOut `json:"data"`
	Object string   `json:"object"`
	Total  int      `json:"total"`
}

// ListParams holds query parameters for listing batch jobs.
type ListParams struct {
	Page         *int
	PageSize     *int
	Model        *string
	AgentID      *string
	CreatedAfter *string
	CreatedByMe  *bool
	Status       []string
	OrderBy      *string
}

// DeleteResponse is the response from deleting a batch job.
type DeleteResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Deleted bool   `json:"deleted"`
}
