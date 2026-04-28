# mistral-go-sdk

Go SDK for Mistral's agentic stack — Workflows, Conversations, Connectors — with full coverage of the rest of the Mistral AI API.

<!-- Badges -->
[![Go Reference](https://pkg.go.dev/badge/github.com/VikingOwl91/mistral-go-sdk.svg)](https://pkg.go.dev/github.com/VikingOwl91/mistral-go-sdk)
![Go Version](https://img.shields.io/badge/go-1.26-blue)
![License](https://img.shields.io/badge/license-MIT-green)

## What this is for

Mistral's `/v1/chat/completions` endpoint is OpenAI-wire-compatible, so if all you need is plain chat or tool calling, pointing any OpenAI Go client at `https://api.mistral.ai/v1` already works. **This SDK exists for the rest of Mistral's surface** — the parts no other Go client covers:

- **Workflows** — durable, long-running executions with signals, queries, updates, and full event streaming.
- **Conversations API** — server-stored multi-turn state with streaming events, agent handoffs, and human-in-the-loop tool confirmations.
- **Connectors / MCP** — manage MCP connectors and bind them to workflow executions.
- **Agents CRUD** — create, version, and alias agents (not just `agents/completions`).
- **Libraries / Documents** — RAG document stores with reprocess, status, and signed-URL retrieval.
- **OCR, Audio (transcription + TTS + voices), Fine-tuning, Batch, Observability, Moderation, Classification.**

If your use case is "Go program calls Mistral and gets a response back," any OpenAI-compatible client is the easier path. If you need agentic workflows or stateful conversations from Go, this is the only option.

## Why this SDK

- **Zero dependencies.** Stdlib-only. No `go.sum`, no transitive tree, no supply chain.
- **Hand-written, not generated.** Idiomatic Go with sealed interfaces, discriminated unions, and functional options — not a Speakeasy/OpenAPI dump with `any` everywhere.
- **Forward-compatible types.** `UnknownEntry`, `UnknownEvent`, `UnknownMessage`, `UnknownChunk`, `UnknownAgentTool` capture raw JSON instead of returning errors — when Mistral ships a new event type, your code keeps running.
- **Typed streaming.** Generic pull-based `Stream[T]` iterator — no channels, no goroutines, no leaks. Just `Next()` / `Current()` / `Err()` / `Close()`.
- **Test-driven.** 297 tests, race-clean. Every endpoint has a mock-server unit test; integration tests run against the real API behind a build tag.

## Install

```sh
go get github.com/VikingOwl91/mistral-go-sdk
```

## Quick Start

### Execute a workflow with connector bindings

The most differentiated thing this SDK does. Run a registered workflow, hand it the MCP connectors it needs, and block until completion:

```go
package main

import (
    "context"
    "encoding/json"
    "fmt"
    "log"

    mistral "github.com/VikingOwl91/mistral-go-sdk"
    "github.com/VikingOwl91/mistral-go-sdk/workflow"
)

func main() {
    client := mistral.NewClient("sk-your-api-key")

    creds := "work-account"
    resp, err := client.ExecuteWorkflowAndWait(context.Background(), "my-workflow", &workflow.ExecutionRequest{
        Input: map[string]any{
            "topic": "Q3 earnings",
        },
        Extensions: workflow.BuildConnectorExtensions(
            workflow.ConnectorSlot{ConnectorName: "gmail"},
            workflow.ConnectorSlot{ConnectorName: "notion", CredentialsName: &creds},
        ),
    })
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("execution %s finished with status %s\n", resp.ExecutionID, resp.Status)
    if resp.Status == workflow.ExecutionCompleted {
        result, _ := json.MarshalIndent(resp.Result, "", "  ")
        fmt.Println(string(result))
    }
}
```

### Stream a conversation with tool-call confirmation

Conversations keep server-side state across turns. Stream events, accept or deny pending tool calls, and continue:

```go
import "github.com/VikingOwl91/mistral-go-sdk/conversation"

stream, err := client.StartConversationStream(ctx, &conversation.StartRequest{
    AgentID: "ag-your-agent-id",
    Inputs:  conversation.TextInputs("Summarize today's incident review and email it to the on-call channel."),
})
if err != nil {
    log.Fatal(err)
}
defer stream.Close()

for stream.Next() {
    event := stream.Current()
    switch e := event.(type) {
    case *conversation.FunctionCallEvent:
        if e.ConfirmationStatus != nil && *e.ConfirmationStatus == conversation.ConfirmationStatusPending {
            // Hand the user a confirmation prompt; reply on the next AppendConversation.
            confirmations := []conversation.ToolCallConfirmation{
                {ToolCallID: e.ToolCallID, Confirmation: string(conversation.ConfirmationAllow)},
            }
            _, _ = client.AppendConversation(ctx, /* conversationID */ "conv-...", &conversation.AppendRequest{
                ToolConfirmations: confirmations,
            })
        }
    case *conversation.MessageOutputEntry:
        // assistant tokens
    }
}
if err := stream.Err(); err != nil {
    log.Fatal(err)
}
```

### Chat completion (when that's all you need)

```go
import "github.com/VikingOwl91/mistral-go-sdk/chat"

resp, err := client.ChatComplete(ctx, &chat.CompletionRequest{
    Model: "mistral-small-latest",
    Messages: []chat.Message{
        &chat.UserMessage{Content: chat.TextContent("What is the capital of France?")},
    },
})
if err != nil {
    log.Fatal(err)
}
fmt.Println(resp.Choices[0].Message.Content)
```

If this is your whole use case, the OpenAI-wire-compatible endpoint is also a valid path — point any `openai-go` client at `https://api.mistral.ai/v1` and it works. The rest of this SDK is for what the OpenAI shape doesn't cover.

## Configuration

```go
client := mistral.NewClient("sk-...",
    mistral.WithBaseURL("https://custom-endpoint.example.com"),
    mistral.WithHTTPClient(customHTTPClient),
    mistral.WithTimeout(30 * time.Second),
    mistral.WithRetry(3, 500*time.Millisecond),
)
```

| Option | Description |
|--------|-------------|
| `WithBaseURL(url)` | Override the API base URL (default: `https://api.mistral.ai`) |
| `WithHTTPClient(client)` | Use a custom `*http.Client` |
| `WithTimeout(d)` | Set HTTP client timeout (default: 120s) |
| `WithRetry(n, delay)` | Enable retry with exponential backoff + jitter on 429/5xx |

## Error Handling

```go
resp, err := client.ChatComplete(ctx, req)
if err != nil {
    var apiErr *mistral.APIError
    if errors.As(err, &apiErr) {
        fmt.Printf("Status: %d, Message: %s\n", apiErr.StatusCode, apiErr.Message)
    }

    if mistral.IsRateLimit(err) {
        // 429 — back off and retry
    }
    if mistral.IsNotFound(err) {
        // 404 — resource not found
    }
    if mistral.IsAuth(err) {
        // 401 — invalid API key
    }
}
```

## API Coverage

169 public methods on `Client`, grouped by domain:

| Domain | Methods |
|--------|---------|
| **Chat** | `ChatComplete`, `ChatCompleteStream` |
| **FIM** | `FIMComplete`, `FIMCompleteStream` |
| **Agents (completions)** | `AgentsComplete`, `AgentsCompleteStream` |
| **Agents (CRUD)** | `CreateAgent`, `ListAgents`, `GetAgent`, `UpdateAgent`, `DeleteAgent`, `UpdateAgentVersion`, `ListAgentVersions`, `GetAgentVersion`, `SetAgentAlias`, `ListAgentAliases`, `DeleteAgentAlias` |
| **Connectors** | `CreateConnector`, `ListConnectors`, `GetConnector`, `UpdateConnector`, `DeleteConnector`, `GetConnectorAuthURL`, `ListConnectorTools`, `CallConnectorTool` |
| **Conversations** | `StartConversation`, `StartConversationStream`, `AppendConversation`, `AppendConversationStream`, `RestartConversation`, `RestartConversationStream`, `GetConversation`, `ListConversations`, `DeleteConversation`, `GetConversationHistory`, `GetConversationMessages` |
| **Models** | `ListModels`, `GetModel`, `DeleteModel` |
| **Files** | `UploadFile`, `ListFiles`, `GetFile`, `DeleteFile`, `GetFileContent`, `GetFileURL` |
| **Embeddings** | `CreateEmbeddings` |
| **Fine-tuning** | `CreateFineTuningJob`, `ListFineTuningJobs`, `GetFineTuningJob`, `CancelFineTuningJob`, `StartFineTuningJob`, `UpdateFineTunedModel`, `ArchiveFineTunedModel`, `UnarchiveFineTunedModel` |
| **Batch** | `CreateBatchJob`, `ListBatchJobs`, `GetBatchJob`, `CancelBatchJob`, `DeleteBatchJob` |
| **OCR** | `OCR` |
| **Audio (transcription)** | `Transcribe`, `TranscribeStream` |
| **Audio (speech)** | `Speech`, `SpeechStream` |
| **Audio (voices)** | `ListVoices`, `CreateVoice`, `GetVoice`, `UpdateVoice`, `DeleteVoice`, `GetVoiceSampleAudio` |
| **Libraries** | `CreateLibrary`, `ListLibraries`, `GetLibrary`, `UpdateLibrary`, `DeleteLibrary`, `UploadDocument`, `ListDocuments`, `GetDocument`, `UpdateDocument`, `DeleteDocument`, `GetDocumentTextContent`, `GetDocumentStatus`, `GetDocumentSignedURL`, `GetDocumentExtractedTextSignedURL`, `ReprocessDocument`, `ListLibrarySharing`, `ShareLibrary`, `UnshareLibrary` |
| **Moderation** | `Moderate`, `ModerateChat` |
| **Classification** | `Classify`, `ClassifyChat` |
| **Observability (campaigns)** | `CreateCampaign`, `ListCampaigns`, `GetCampaign`, `DeleteCampaign`, `GetCampaignStatus`, `ListCampaignEvents` |
| **Observability (events)** | `SearchChatCompletionEvents`, `SearchChatCompletionEventIDs`, `GetChatCompletionEvent`, `GetSimilarChatCompletionEvents`, `JudgeChatCompletionEvent` |
| **Observability (fields)** | `GetChatCompletionFields`, `GetChatCompletionFieldOptions`, `GetChatCompletionFieldOptionsCounts` |
| **Observability (judges)** | `CreateJudge`, `ListJudges`, `GetJudge`, `UpdateJudge`, `DeleteJudge`, `JudgeConversation` |
| **Observability (datasets)** | `CreateDataset`, `ListDatasets`, `GetDataset`, `UpdateDataset`, `DeleteDataset`, `ExportDatasetToJSONL`, `ListDatasetRecords`, `CreateDatasetRecord`, `GetDatasetRecord`, `UpdateDatasetRecordPayload`, `UpdateDatasetRecordProperties`, `DeleteDatasetRecord`, `BulkDeleteDatasetRecords`, `JudgeDatasetRecord`, `ImportDatasetFromCampaign`, `ImportDatasetFromExplorer`, `ImportDatasetFromFile`, `ImportDatasetFromPlayground`, `ImportDatasetFromDataset`, `ListDatasetTasks`, `GetDatasetTask` |
| **Workflows (CRUD)** | `ListWorkflows`, `GetWorkflow`, `UpdateWorkflow`, `ArchiveWorkflow`, `UnarchiveWorkflow`, `ExecuteWorkflow`, `ExecuteWorkflowAndWait` |
| **Workflows (registrations)** | `ListWorkflowRegistrations`, `GetWorkflowRegistration`, `ExecuteWorkflowRegistration` |
| **Workflows (executions)** | `GetWorkflowExecution`, `GetWorkflowExecutionHistory`, `StreamWorkflowExecution`, `SignalWorkflowExecution`, `QueryWorkflowExecution`, `UpdateWorkflowExecution`, `TerminateWorkflowExecution`, `CancelWorkflowExecution`, `ResetWorkflowExecution`, `BatchCancelWorkflowExecutions`, `BatchTerminateWorkflowExecutions` |
| **Workflows (trace)** | `GetWorkflowExecutionTraceOTel`, `GetWorkflowExecutionTraceSummary`, `GetWorkflowExecutionTraceEvents` |
| **Workflows (events)** | `StreamWorkflowEvents`, `ListWorkflowEvents` |
| **Workflows (deployments)** | `ListWorkflowDeployments`, `GetWorkflowDeployment` |
| **Workflows (metrics)** | `GetWorkflowMetrics` |
| **Workflows (runs)** | `ListWorkflowRuns`, `GetWorkflowRun`, `GetWorkflowRunHistory` |
| **Workflows (schedules)** | `ListWorkflowSchedules`, `ScheduleWorkflow`, `UnscheduleWorkflow` |
| **Workflows (workers)** | `GetWorkflowWorkerInfo` |

## Comparison

There is no official Go SDK from Mistral AI (only Python and TypeScript). The main community options:

| Feature | mistral-go-sdk | [Gage-Technologies](https://github.com/Gage-Technologies/mistral-go) | [robertjkeck2](https://github.com/robertjkeck2/mistral-go) | [AuxData-ai](https://github.com/AuxData-ai/mistral-go) |
|---------|:-:|:-:|:-:|:-:|
| Chat / Streaming | Yes | Yes | Yes | Yes |
| FIM | Yes | Yes | No | Yes |
| Embeddings | Yes | Yes | Yes | Yes |
| Tool calling | Yes | No | No | No |
| Agents (completions + CRUD) | Yes | No | No | No |
| Connectors (MCP) | Yes | No | No | No |
| Conversations API | Yes | No | No | No |
| Libraries / Documents | Yes | No | No | No |
| Fine-tuning / Batch | Yes | No | No | No |
| OCR | Yes | No | No | Yes |
| Audio (transcription + TTS + voices) | Yes | No | No | No |
| Workflows API | Yes | No | No | No |
| Observability (beta) | Yes | No | No | No |
| Moderation / Classification | Yes | No | No | No |
| Vision (multimodal) | Yes | No | No | Yes |
| Zero dependencies | Yes | test-only (testify) | test-only (testify) | test-only (testify) |
| Forward-compatible types | Yes | No | No | No |
| Last updated | 2026 | Jun 2024 | Jan 2024 | ~2025 (fork of Gage) |

## Upstream Reference

This SDK tracks the [official Mistral OpenAPI spec](https://github.com/mistralai/platform-docs-public/blob/main/openapi.yaml) as its primary reference for API surface and type definitions. A daily GitHub Action monitors the spec for changes and refreshes a tracking issue when updates are detected.

The [Mistral Python SDK](https://github.com/mistralai/client-python) is used as a secondary reference for implementation patterns.

| SDK Version | Upstream Python SDK | Upstream OpenAPI |
|-------------|---------------------|------------------|
| v1.4.0 | v2.4.3 (excl. RAG ingestion-pipeline beta) | v1.0.0 |
| v1.3.0 | v2.3.0 | v0.1.104 |
| v1.2.1 | v2.2.0 | — |
| v1.2.0 | v2.2.0 | — |
| v1.1.0 | v2.1.3 | — |
| v1.0.0 | v2.0.4 | — |

## License

[MIT](LICENSE)
