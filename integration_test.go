//go:build integration

package mistral

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/VikingOwl91/mistral-go-sdk/chat"
	"github.com/VikingOwl91/mistral-go-sdk/embedding"
)

func integrationClient(t *testing.T) *Client {
	t.Helper()
	key := os.Getenv("MISTRAL_API_KEY")
	if key == "" {
		t.Skip("MISTRAL_API_KEY not set")
	}
	return NewClient(key)
}

func TestIntegration_ListModels(t *testing.T) {
	client := integrationClient(t)

	resp, err := client.ListModels(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.Data) == 0 {
		t.Fatal("expected at least one model")
	}
	if resp.Data[0].ID == "" {
		t.Error("expected model to have an ID")
	}
}

func TestIntegration_ChatComplete(t *testing.T) {
	client := integrationClient(t)

	resp, err := client.ChatComplete(context.Background(), &chat.CompletionRequest{
		Model: "mistral-small-latest",
		Messages: []chat.Message{
			&chat.UserMessage{Content: chat.TextContent("Reply with exactly: hello")},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	if len(resp.Choices) == 0 {
		t.Fatal("no choices returned")
	}
	content := resp.Choices[0].Message.Content.String()
	if !strings.Contains(strings.ToLower(content), "hello") {
		t.Errorf("expected 'hello' in response, got %q", content)
	}
}

func TestIntegration_ChatCompleteStream(t *testing.T) {
	client := integrationClient(t)

	stream, err := client.ChatCompleteStream(context.Background(), &chat.CompletionRequest{
		Model: "mistral-small-latest",
		Messages: []chat.Message{
			&chat.UserMessage{Content: chat.TextContent("Say hi")},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer stream.Close()

	chunks := 0
	for stream.Next() {
		chunks++
	}
	if err := stream.Err(); err != nil {
		t.Fatal(err)
	}
	if chunks == 0 {
		t.Error("expected at least one stream chunk")
	}
}

func TestIntegration_Embeddings(t *testing.T) {
	client := integrationClient(t)

	resp, err := client.CreateEmbeddings(context.Background(), &embedding.Request{
		Model: "mistral-embed",
		Input: []string{"The quick brown fox"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.Data) != 1 {
		t.Fatalf("expected 1 embedding, got %d", len(resp.Data))
	}
	if len(resp.Data[0].Embedding) == 0 {
		t.Error("embedding vector is empty")
	}
}

func TestIntegration_ErrorHandling(t *testing.T) {
	client := integrationClient(t)

	_, err := client.GetModel(context.Background(), "definitely-not-a-real-model-id")
	if err == nil {
		t.Fatal("expected error for nonexistent model")
	}
	if !IsNotFound(err) {
		t.Logf("got non-404 error (may vary by API): %v", err)
	}
}
