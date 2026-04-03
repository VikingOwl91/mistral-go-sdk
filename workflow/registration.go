package workflow

// Registration represents a workflow registration.
type Registration struct {
	ID                         string          `json:"id"`
	WorkflowID                 string          `json:"workflow_id"`
	Definition                 *CodeDefinition `json:"definition,omitempty"`
	DeploymentID               *string         `json:"deployment_id,omitempty"`
	CompatibleWithChatAssistant bool           `json:"compatible_with_chat_assistant,omitempty"`
	// Deprecated: use DeploymentID instead. Will be removed in a future release.
	TaskQueue  string    `json:"task_queue"`
	Workflow   *Workflow `json:"workflow,omitempty"`
	CreatedAt  string    `json:"created_at"`
	UpdatedAt  string    `json:"updated_at"`
}

// RegistrationListResponse is the response from listing workflow registrations.
type RegistrationListResponse struct {
	Registrations []Registration `json:"registrations"`
	NextCursor    *string        `json:"next_cursor,omitempty"`
}

// RegistrationListParams holds query parameters for listing registrations.
type RegistrationListParams struct {
	WorkflowID               *string
	TaskQueue                *string
	ActiveOnly               *bool
	IncludeShared            *bool
	WorkflowSearch           *string
	Archived                 *bool
	WithWorkflow             *bool
	AvailableInChatAssistant *bool
	Limit                    *int
	Cursor                   *string
}

// RegistrationGetParams holds query parameters for getting a registration.
type RegistrationGetParams struct {
	WithWorkflow  *bool
	IncludeShared *bool
}
