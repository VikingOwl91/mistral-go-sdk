// Package agents provides types for the Mistral agents API,
// including agent CRUD operations and agent chat completions.
//
// # Tool Types
//
// Agents support multiple tool types via the [AgentTool] sealed interface:
// [FunctionTool], [WebSearchTool], [WebSearchPremiumTool],
// [CodeInterpreterTool], [ImageGenerationTool], [DocumentLibraryTool],
// and [ConnectorTool] for custom connectors.
package agents
