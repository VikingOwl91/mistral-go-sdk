package mistral_test

import (
	"context"
	"fmt"
	"log"

	mistral "github.com/VikingOwl91/mistral-go-sdk"
	"github.com/VikingOwl91/mistral-go-sdk/chat"
	"github.com/VikingOwl91/mistral-go-sdk/embedding"
)

func ExampleNewClient() {
	client := mistral.NewClient("sk-your-api-key")
	_ = client // use client to call API methods
}

func ExampleClient_ChatComplete() {
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

func ExampleClient_ChatCompleteStream() {
	client := mistral.NewClient("sk-your-api-key")

	stream, err := client.ChatCompleteStream(context.Background(), &chat.CompletionRequest{
		Model: "mistral-small-latest",
		Messages: []chat.Message{
			&chat.UserMessage{Content: chat.TextContent("Tell me a short joke.")},
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
}

func ExampleClient_CreateEmbeddings() {
	client := mistral.NewClient("sk-your-api-key")

	resp, err := client.CreateEmbeddings(context.Background(), &embedding.Request{
		Model: "mistral-embed",
		Input: []string{"Hello world"},
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Embedding dimension: %d\n", len(resp.Data[0].Embedding))
}

func ExampleIsRateLimit() {
	client := mistral.NewClient("sk-your-api-key")

	_, err := client.ChatComplete(context.Background(), &chat.CompletionRequest{
		Model: "mistral-small-latest",
		Messages: []chat.Message{
			&chat.UserMessage{Content: chat.TextContent("Hi")},
		},
	})
	if mistral.IsRateLimit(err) {
		fmt.Println("Rate limited — back off and retry")
	}
}

func ExampleIsNotFound() {
	client := mistral.NewClient("sk-your-api-key")

	_, err := client.GetModel(context.Background(), "nonexistent-model")
	if mistral.IsNotFound(err) {
		fmt.Println("Model not found")
	}
}
