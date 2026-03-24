// Package mistral provides an idiomatic Go client for the Mistral AI API.
//
// Create a client with your API key, then call methods for each endpoint:
//
//	client := mistral.NewClient("sk-...")
//
//	// Chat completion
//	resp, err := client.ChatComplete(ctx, &chat.CompletionRequest{
//	    Model:    "mistral-small-latest",
//	    Messages: []chat.Message{&chat.UserMessage{Content: chat.TextContent("Hello!")}},
//	})
//
//	// Streaming
//	stream, err := client.ChatCompleteStream(ctx, req)
//	defer stream.Close()
//	for stream.Next() {
//	    chunk := stream.Current()
//	    fmt.Print(chunk.Choices[0].Delta.Content)
//	}
//
// # Configuration
//
// Use functional options to configure the client:
//
//	client := mistral.NewClient("sk-...",
//	    mistral.WithTimeout(30 * time.Second),
//	    mistral.WithRetry(3, 500*time.Millisecond),
//	)
//
// # Error Handling
//
// API errors are returned as *[APIError] values. Use sentinel checkers
// for common cases:
//
//	if mistral.IsRateLimit(err) {
//	    // back off and retry
//	}
//
// # Sub-packages
//
// Types are organized into sub-packages by domain: [chat], [agents],
// [connector], [conversation], [embedding], [model], [file], [finetune],
// [batch], [ocr], [audio], [library], [moderation], [classification],
// [fim], and [observability]. All service methods live directly on [Client].
//
// # Reference
//
// This SDK tracks the official Mistral Python SDK
// (https://github.com/mistralai/client-python) as its upstream reference
// for API surface and type definitions.
package mistral
