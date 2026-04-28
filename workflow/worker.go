package workflow

// WorkerInfo describes the worker scheduler the SDK is connected to.
//
// Returned by GET /v1/workflows/workers/whoami. Useful when running custom
// workers that need to know which scheduler / namespace to connect to.
// For managed deployments, prefer Registration.DeploymentID.
type WorkerInfo struct {
	SchedulerURL string `json:"scheduler_url"`
	Namespace    string `json:"namespace"`
	TLS          bool   `json:"tls"`
}
