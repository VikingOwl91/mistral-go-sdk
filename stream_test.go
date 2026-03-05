package mistral

import (
	"io"
	"strings"
	"testing"
)

type testChunk struct {
	ID      string `json:"id"`
	Content string `json:"content"`
}

func newTestStream(sse string) *Stream[testChunk] {
	body := io.NopCloser(strings.NewReader(sse))
	return newStream[testChunk](body)
}

func TestStream_SingleChunk(t *testing.T) {
	input := "data: {\"id\":\"1\",\"content\":\"hello\"}\n\ndata: [DONE]\n\n"
	s := newTestStream(input)
	defer s.Close()

	if !s.Next() {
		t.Fatalf("expected Next() to return true, err: %v", s.Err())
	}
	chunk := s.Current()
	if chunk.ID != "1" || chunk.Content != "hello" {
		t.Errorf("got %+v", chunk)
	}
	if s.Next() {
		t.Error("expected Next() to return false after [DONE]")
	}
	if s.Err() != nil {
		t.Errorf("unexpected error: %v", s.Err())
	}
}

func TestStream_MultipleChunks(t *testing.T) {
	input := "data: {\"id\":\"1\",\"content\":\"a\"}\n\ndata: {\"id\":\"2\",\"content\":\"b\"}\n\ndata: {\"id\":\"3\",\"content\":\"c\"}\n\ndata: [DONE]\n\n"
	s := newTestStream(input)
	defer s.Close()

	var chunks []testChunk
	for s.Next() {
		chunks = append(chunks, s.Current())
	}
	if s.Err() != nil {
		t.Fatal(s.Err())
	}
	if len(chunks) != 3 {
		t.Fatalf("got %d chunks, want 3", len(chunks))
	}
	if chunks[0].Content != "a" || chunks[1].Content != "b" || chunks[2].Content != "c" {
		t.Errorf("got %+v", chunks)
	}
}

func TestStream_EmptyStream(t *testing.T) {
	s := newTestStream("data: [DONE]\n\n")
	defer s.Close()

	if s.Next() {
		t.Error("expected no chunks before [DONE]")
	}
	if s.Err() != nil {
		t.Errorf("unexpected error: %v", s.Err())
	}
}

func TestStream_InvalidJSON(t *testing.T) {
	input := "data: not-json\n\n"
	s := newTestStream(input)
	defer s.Close()

	if s.Next() {
		t.Error("expected Next() to return false for invalid JSON")
	}
	if s.Err() == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestStream_NextAfterDone(t *testing.T) {
	input := "data: {\"id\":\"1\",\"content\":\"x\"}\n\ndata: [DONE]\n\n"
	s := newTestStream(input)
	defer s.Close()

	s.Next() // consume first chunk
	s.Next() // hits [DONE]

	// Calling Next() again should still return false
	if s.Next() {
		t.Error("expected false after stream is done")
	}
}

func TestStream_NextAfterError(t *testing.T) {
	input := "data: bad\n\n"
	s := newTestStream(input)
	defer s.Close()

	s.Next() // triggers error

	// Calling Next() again should still return false
	if s.Next() {
		t.Error("expected false after error")
	}
}

func TestStream_WithComments(t *testing.T) {
	input := ": keep-alive\ndata: {\"id\":\"1\",\"content\":\"ok\"}\n\n: ping\ndata: [DONE]\n\n"
	s := newTestStream(input)
	defer s.Close()

	if !s.Next() {
		t.Fatalf("expected chunk, err: %v", s.Err())
	}
	if s.Current().Content != "ok" {
		t.Errorf("got %q", s.Current().Content)
	}
	if s.Next() {
		t.Error("expected done after [DONE]")
	}
}

func TestStream_EOFWithoutDone(t *testing.T) {
	input := "data: {\"id\":\"1\",\"content\":\"x\"}\n\n"
	s := newTestStream(input)
	defer s.Close()

	if !s.Next() {
		t.Fatalf("expected chunk, err: %v", s.Err())
	}
	if s.Next() {
		t.Error("expected false at EOF")
	}
	if s.Err() != nil {
		t.Errorf("expected no error at clean EOF, got: %v", s.Err())
	}
}
