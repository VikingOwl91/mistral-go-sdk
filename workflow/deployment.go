package workflow

// Deployment represents a workflow deployment.
type Deployment struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	IsActive  bool   `json:"is_active"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// DeploymentListResponse is the response from listing deployments.
type DeploymentListResponse struct {
	Deployments []Deployment `json:"deployments"`
}

// DeploymentListParams holds query parameters for listing deployments.
type DeploymentListParams struct {
	ActiveOnly   *bool
	WorkflowName *string
}
