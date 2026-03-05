# mistral-go-sdk

An idiomatic Go client for the [Mistral AI API](https://docs.mistral.ai/).

<!-- Badges -->
![Go Version](https://img.shields.io/badge/go-1.26-blue)
![License](https://img.shields.io/badge/license-MIT-green)

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

75 public methods on `Client`, grouped by domain:

| Domain | Methods |
|--------|---------|
| **Chat** | `ChatComplete`, `ChatCompleteStream` |
| **FIM** | `FIMComplete`, `FIMCompleteStream` |
| **Agents (completions)** | `AgentsComplete`, `AgentsCompleteStream` |
| **Agents (CRUD)** | `CreateAgent`, `ListAgents`, `GetAgent`, `UpdateAgent`, `DeleteAgent`, `UpdateAgentVersion`, `ListAgentVersions`, `GetAgentVersion`, `SetAgentAlias`, `ListAgentAliases`, `DeleteAgentAlias` |
| **Conversations** | `StartConversation`, `StartConversationStream`, `AppendConversation`, `AppendConversationStream`, `RestartConversation`, `RestartConversationStream`, `GetConversation`, `ListConversations`, `DeleteConversation`, `GetConversationHistory`, `GetConversationMessages` |
| **Models** | `ListModels`, `GetModel`, `DeleteModel` |
| **Files** | `UploadFile`, `ListFiles`, `GetFile`, `DeleteFile`, `GetFileContent`, `GetFileURL` |
| **Embeddings** | `CreateEmbeddings` |
| **Fine-tuning** | `CreateFineTuningJob`, `ListFineTuningJobs`, `GetFineTuningJob`, `CancelFineTuningJob`, `StartFineTuningJob`, `UpdateFineTunedModel`, `ArchiveFineTunedModel`, `UnarchiveFineTunedModel` |
| **Batch** | `CreateBatchJob`, `ListBatchJobs`, `GetBatchJob`, `CancelBatchJob` |
| **OCR** | `OCR` |
| **Audio** | `Transcribe`, `TranscribeStream` |
| **Libraries** | `CreateLibrary`, `ListLibraries`, `GetLibrary`, `UpdateLibrary`, `DeleteLibrary`, `UploadDocument`, `ListDocuments`, `GetDocument`, `UpdateDocument`, `DeleteDocument`, `GetDocumentTextContent`, `GetDocumentStatus`, `GetDocumentSignedURL`, `GetDocumentExtractedTextSignedURL`, `ReprocessDocument`, `ListLibrarySharing`, `ShareLibrary`, `UnshareLibrary` |
| **Moderation** | `Moderate`, `ModerateChat` |
| **Classification** | `Classify`, `ClassifyChat` |

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

## License

[MIT](LICENSE)
