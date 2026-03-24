# mistral-go-sdk

The most complete Go client for the [Mistral AI API](https://docs.mistral.ai/).

<!-- Badges -->
[![Go Reference](https://pkg.go.dev/badge/somegit.dev/vikingowl/mistral-go-sdk.svg)](https://pkg.go.dev/somegit.dev/vikingowl/mistral-go-sdk)
![Go Version](https://img.shields.io/badge/go-1.26-blue)
![License](https://img.shields.io/badge/license-MIT-green)

## Why This SDK?

**Zero dependencies.** The entire SDK — including tests — uses only the Go standard library. No `go.sum`, no transitive dependency tree to audit, no version conflicts, no supply chain risk.

**Full API coverage.** 116 methods across every Mistral endpoint — including Connectors, Audio Speech/Voices, Conversations, Agents CRUD, Libraries, OCR, Observability, Fine-tuning, and Batch Jobs. No other Go SDK covers Conversations, Connectors, or Observability.

**Typed streaming.** A generic pull-based `Stream[T]` iterator — no channels, no goroutines, no leaks. Just `Next()` / `Current()` / `Err()` / `Close()`.

**Forward-compatible.** Unknown types (`UnknownEntry`, `UnknownEvent`, `UnknownMessage`, `UnknownChunk`, `UnknownAgentTool`) capture raw JSON instead of returning errors. When Mistral ships a new message role or event type, your code keeps running — it doesn't panic.

**Hand-written, not generated.** Idiomatic Go with sealed interfaces, discriminated unions, and functional options — not a Speakeasy/OpenAPI auto-gen dump with `any` everywhere.

**Test-driven.** 193 tests with race detection clean. Every endpoint tested against mock servers; integration tests against the real API.

## Install

```sh
go get somegit.dev/vikingowl/mistral-go-sdk
```

## Quick Start

### Chat Completion

```go
package main

import (
    "context"
    "fmt"
    "log"

    mistral "somegit.dev/vikingowl/mistral-go-sdk"
    "somegit.dev/vikingowl/mistral-go-sdk/chat"
)

func main() {
    client := mistral.NewClient("sk-your-api-key")

    resp, err := client.ChatComplete(context.Background(), &chat.CompletionRequest{
        Model: "mistral-small-latest",
        Messages: []chat.Message{
            &chat.UserMessage{Content: chat.TextContent("What is the capital of France?")},
        },
    })
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(resp.Choices[0].Message.Content)
}
```

### Streaming

```go
stream, err := client.ChatCompleteStream(ctx, &chat.CompletionRequest{
    Model: "mistral-small-latest",
    Messages: []chat.Message{
        &chat.UserMessage{Content: chat.TextContent("Tell me a joke.")},
    },
})
if err != nil {
    log.Fatal(err)
}
defer stream.Close()

for stream.Next() {
    chunk := stream.Current()
    if len(chunk.Choices) > 0 {
        fmt.Print(chunk.Choices[0].Delta.Content)
    }
}
if err := stream.Err(); err != nil {
    log.Fatal(err)
}
```

### Tool Calling

```go
resp, err := client.ChatComplete(ctx, &chat.CompletionRequest{
    Model: "mistral-small-latest",
    Messages: []chat.Message{
        &chat.UserMessage{Content: chat.TextContent("What's the weather in Paris?")},
    },
    Tools: []chat.Tool{{
        Type: "function",
        Function: chat.Function{
            Name:        "get_weather",
            Description: "Get weather for a city",
            Parameters: map[string]any{
                "type": "object",
                "properties": map[string]any{
                    "city": map[string]any{"type": "string"},
                },
                "required": []string{"city"},
            },
        },
    }},
})
```

### Conversations

```go
import "somegit.dev/vikingowl/mistral-go-sdk/conversation"

resp, err := client.StartConversation(ctx, &conversation.StartRequest{
    AgentID: "ag-your-agent-id",
    Inputs:  conversation.TextInputs("Hello, agent!"),
})

// Stream events
stream, err := client.AppendConversationStream(ctx, resp.ConversationID, &conversation.AppendRequest{
    Inputs: conversation.TextInputs("Follow-up question"),
})
defer stream.Close()
for stream.Next() {
    event := stream.Current()
    // handle typed events
}
```

## API Coverage

116 public methods on `Client`, grouped by domain:

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
| **Batch** | `CreateBatchJob`, `ListBatchJobs`, `GetBatchJob`, `CancelBatchJob` |
| **OCR** | `OCR` |
| **Audio (transcription)** | `Transcribe`, `TranscribeStream` |
| **Audio (speech)** | `Speech`, `SpeechStream` |
| **Audio (voices)** | `ListVoices`, `CreateVoice`, `GetVoice`, `UpdateVoice`, `DeleteVoice`, `GetVoiceSampleAudio` |
| **Libraries** | `CreateLibrary`, `ListLibraries`, `GetLibrary`, `UpdateLibrary`, `DeleteLibrary`, `UploadDocument`, `ListDocuments`, `GetDocument`, `UpdateDocument`, `DeleteDocument`, `GetDocumentTextContent`, `GetDocumentStatus`, `GetDocumentSignedURL`, `GetDocumentExtractedTextSignedURL`, `ReprocessDocument`, `ListLibrarySharing`, `ShareLibrary`, `UnshareLibrary` |
| **Moderation** | `Moderate`, `ModerateChat` |
| **Classification** | `Classify`, `ClassifyChat` |
| **Observability (campaigns)** | `CreateCampaign`, `ListCampaigns`, `GetCampaign`, `DeleteCampaign`, `GetCampaignStatus`, `ListCampaignEvents` |
| **Observability (events)** | `SearchChatCompletionEvents`, `SearchChatCompletionEventIDs`, `GetChatCompletionEvent`, `GetSimilarChatCompletionEvents`, `JudgeChatCompletionEvent` |
| **Observability (judges)** | `CreateJudge`, `ListJudges`, `GetJudge`, `UpdateJudge`, `DeleteJudge`, `JudgeConversation` |
| **Observability (datasets)** | `CreateDataset`, `ListDatasets`, `GetDataset`, `UpdateDataset`, `DeleteDataset`, `ExportDatasetToJSONL`, `ListDatasetRecords`, `CreateDatasetRecord`, `GetDatasetRecord`, `UpdateDatasetRecordPayload`, `UpdateDatasetRecordProperties`, `DeleteDatasetRecord`, `BulkDeleteDatasetRecords`, `JudgeDatasetRecord`, `ImportDatasetFromCampaign`, `ImportDatasetFromExplorer`, `ImportDatasetFromFile`, `ImportDatasetFromPlayground`, `ImportDatasetFromDataset`, `ListDatasetTasks`, `GetDatasetTask` |

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
| Observability (beta) | Yes | No | No | No |
| Moderation / Classification | Yes | No | No | No |
| Vision (multimodal) | Yes | No | No | Yes |
| Zero dependencies | Yes | test-only (testify) | test-only (testify) | test-only (testify) |
| Forward-compatible types | Yes | No | No | No |
| Last updated | 2026 | Jun 2024 | Jan 2024 | ~2025 (fork of Gage) |

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

## Upstream Reference

This SDK tracks the [official Mistral Python SDK](https://github.com/mistralai/client-python)
as its upstream reference for API surface and type definitions.

| SDK Version | Upstream Python SDK |
|-------------|---------------------|
| v1.1.0 | v2.1.3 |
| v1.0.0 | v2.0.4 |

## License

[MIT](LICENSE)
