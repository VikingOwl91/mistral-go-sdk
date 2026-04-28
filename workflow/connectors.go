package workflow

// ConnectorSlot declares a connector dependency for a workflow execution.
//
// Pass a slice of slots to BuildConnectorExtensions to produce the
// nested map expected on ExecutionRequest.Extensions.
type ConnectorSlot struct {
	ConnectorName   string  `json:"connector_name"`
	CredentialsName *string `json:"credentials_name,omitempty"`
}

// ConnectorBindings is the bindings list inside ConnectorExtensions.
type ConnectorBindings struct {
	Bindings []ConnectorSlot `json:"bindings"`
}

// ConnectorExtensions is the value of the "mistralai" key in workflow extensions.
type ConnectorExtensions struct {
	Connectors ConnectorBindings `json:"connectors"`
}

// WorkflowExtensions is the top-level shape of the extensions field
// expected by the workflow execute endpoint when binding connectors.
type WorkflowExtensions struct {
	Mistralai ConnectorExtensions `json:"mistralai"`
}

// BuildConnectorExtensions returns the value to set on
// ExecutionRequest.Extensions for the given connector slots.
//
// The result is a map[string]any so callers can merge in additional
// extension keys without colliding with the connector wire shape.
func BuildConnectorExtensions(slots ...ConnectorSlot) map[string]any {
	return map[string]any{
		"mistralai": ConnectorExtensions{
			Connectors: ConnectorBindings{Bindings: slots},
		},
	}
}

// ConnectorAuthStatus is the state of an OAuth flow emitted by a
// connector-auth custom task event.
type ConnectorAuthStatus string

const (
	ConnectorAuthWaitingForAuth ConnectorAuthStatus = "waiting_for_auth"
	ConnectorAuthConnected      ConnectorAuthStatus = "connected"
	ConnectorAuthAccessDenied   ConnectorAuthStatus = "access_denied"
	ConnectorAuthTimedOut       ConnectorAuthStatus = "timed_out"
	ConnectorAuthError          ConnectorAuthStatus = "error"
)

// ConnectorAuthTaskState is the payload of a custom task event of type
// "connector-auth", emitted while a workflow waits for OAuth completion.
type ConnectorAuthTaskState struct {
	ConnectorName string              `json:"connector_name"`
	ConnectorID   string              `json:"connector_id"`
	Status        ConnectorAuthStatus `json:"status"`
	AuthURL       *string             `json:"auth_url,omitempty"`
	Message       *string             `json:"message,omitempty"`
}
