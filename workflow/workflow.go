package workflow

// Workflow represents a workflow definition.
type Workflow struct {
	ID                       string  `json:"id"`
	Name                     string  `json:"name"`
	DisplayName              *string `json:"display_name,omitempty"`
	Description              *string `json:"description,omitempty"`
	OwnerID                  string  `json:"owner_id"`
	WorkspaceID              string  `json:"workspace_id"`
	AvailableInChatAssistant bool    `json:"available_in_chat_assistant"`
	Archived                 bool    `json:"archived"`
	CreatedAt                string  `json:"created_at"`
	UpdatedAt                string  `json:"updated_at"`
}

// WorkflowUpdateRequest is the request body for updating a workflow.
type WorkflowUpdateRequest struct {
	DisplayName              *string `json:"display_name,omitempty"`
	Description              *string `json:"description,omitempty"`
	AvailableInChatAssistant *bool   `json:"available_in_chat_assistant,omitempty"`
}

// WorkflowListResponse is the response from listing workflows.
type WorkflowListResponse struct {
	Workflows  []Workflow `json:"workflows"`
	NextCursor *string    `json:"next_cursor,omitempty"`
}

// WorkflowListParams holds query parameters for listing workflows.
type WorkflowListParams struct {
	ActiveOnly               *bool
	IncludeShared            *bool
	AvailableInChatAssistant *bool
	Archived                 *bool
	Cursor                   *string
	Limit                    *int
}

// WorkflowArchiveResponse is the response from archiving/unarchiving a workflow.
type WorkflowArchiveResponse struct {
	ID       string `json:"id"`
	Archived bool   `json:"archived"`
}
