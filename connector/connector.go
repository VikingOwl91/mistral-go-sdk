package connector

import "encoding/json"

// Visibility controls who can see a connector or tool.
type Visibility string

const (
	VisibilitySharedGlobal    Visibility = "shared_global"
	VisibilitySharedOrg       Visibility = "shared_org"
	VisibilitySharedWorkspace Visibility = "shared_workspace"
	VisibilityPrivate         Visibility = "private"
)

// AuthData holds OAuth2 client credentials for a connector.
type AuthData struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

// Connector represents a registered MCP connector.
type Connector struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	CreatedAt   string  `json:"created_at"`
	ModifiedAt  string  `json:"modified_at"`
	Server      *string `json:"server,omitempty"`
	AuthType    *string `json:"auth_type,omitempty"`
}

// CreateRequest creates a new connector.
type CreateRequest struct {
	Name         string            `json:"name"`
	Description  string            `json:"description"`
	Server       string            `json:"server"`
	IconURL      *string           `json:"icon_url,omitempty"`
	Visibility   *Visibility       `json:"visibility,omitempty"`
	Headers      map[string]string `json:"headers,omitempty"`
	AuthData     *AuthData         `json:"auth_data,omitempty"`
	SystemPrompt *string           `json:"system_prompt,omitempty"`
}

// UpdateRequest updates an existing connector.
type UpdateRequest struct {
	Name              *string           `json:"name,omitempty"`
	Description       *string           `json:"description,omitempty"`
	IconURL           *string           `json:"icon_url,omitempty"`
	SystemPrompt      *string           `json:"system_prompt,omitempty"`
	Server            *string           `json:"server,omitempty"`
	Headers           map[string]string `json:"headers,omitempty"`
	AuthData          *AuthData         `json:"auth_data,omitempty"`
	ConnectionConfig  map[string]any    `json:"connection_config,omitempty"`
	ConnectionSecrets map[string]any    `json:"connection_secrets,omitempty"`
}

// AuthURLResponse is the response from getting a connector's OAuth URL.
type AuthURLResponse struct {
	AuthURL string `json:"auth_url"`
	TTL     int    `json:"ttl"`
}

// CallToolRequest is the request body for calling a connector tool.
type CallToolRequest struct {
	Arguments map[string]any `json:"arguments,omitempty"`
}

// CallToolResponse is the response from calling a connector tool.
// Content is left as raw JSON because the upstream API returns a union
// of 5 content types (text, image, audio, resource link, embedded resource).
type CallToolResponse struct {
	Content  json.RawMessage `json:"content"`
	Metadata map[string]any  `json:"metadata,omitempty"`
}

// Tool represents a tool exposed by a connector.
type Tool struct {
	ID           string         `json:"id"`
	Name         string         `json:"name"`
	Description  *string        `json:"description,omitempty"`
	Visibility   Visibility     `json:"visibility,omitempty"`
	CreatedAt    string         `json:"created_at,omitempty"`
	ModifiedAt   string         `json:"modified_at,omitempty"`
	SystemPrompt *string        `json:"system_prompt,omitempty"`
	JsonSchema   map[string]any `json:"jsonschema,omitempty"`
	Active       *bool          `json:"active,omitempty"`
}

// ListParams holds query parameters for listing connectors.
type ListParams struct {
	Page     *int
	PageSize *int
}

// ListToolsParams holds query parameters for listing connector tools.
type ListToolsParams struct {
	Page     *int
	PageSize *int
	Refresh  *bool
}
