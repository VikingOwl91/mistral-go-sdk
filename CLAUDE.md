# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project

Idiomatic Go SDK for the Mistral AI API. Module path: `somegit.dev/vikingowl/mistral-go-sdk`. Requires Go 1.26+. Zero external dependencies (stdlib only). Tracks the upstream [Mistral Python SDK](https://github.com/mistralai/client-python) as reference for API surface and type definitions.

## Repository layout

- **Working directory**: `mistral-go-sdk/` — the Go SDK source. All development happens here.
- **`../client-python/`**: Clone of the upstream Mistral Python SDK. Read-only reference — pull/update it when checking for upstream API changes, but never modify it.

## Commands

```bash
# Run all unit tests
go test ./...

# Run a single test
go test -run TestChatComplete_Success

# Run integration tests (requires MISTRAL_API_KEY env var)
go test -tags=integration ./...

# Vet and build
go vet ./...
go build ./...
```

No Makefile, linter config, or code generation tooling — standard `go test` / `go vet` / `go build`.

## Architecture

### Two-layer design: types in sub-packages, methods on `*Client`

Sub-packages (`chat/`, `agents/`, `conversation/`, `embedding/`, `model/`, `file/`, `finetune/`, `batch/`, `ocr/`, `audio/`, `library/`, `moderation/`, `classification/`, `fim/`) are **types-only** — they define request/response structs and enums but contain no HTTP logic. All service methods live on `*Client` in the root package, prefix-namespaced by domain (e.g. `ChatComplete`, `AgentsComplete`, `CreateFineTuningJob`, `UploadFile`).

### HTTP internals (request.go)

All HTTP flows route through a small set of unexported helpers on `*Client`:
- `do()` — raw HTTP with auth headers + retry
- `doJSON()` — JSON marshal request → `do()` → unmarshal response
- `doStream()` — JSON request → raw `*http.Response` for SSE
- `doMultipart()` / `doMultipartStream()` — multipart file upload variants
- `doRetry()` — retry loop with exponential backoff + jitter + `Retry-After` parsing

### Streaming

Generic `Stream[T]` type wraps SSE (`sseReader`) with `Next()`/`Current()`/`Err()`/`Close()` iterator pattern. Typed wrappers `EventStream` (conversations) and `AudioStream` (transcription) unmarshal `json.RawMessage` into domain-specific event types.

### Sealed interfaces for discriminated unions

Polymorphic API types use **sealed interfaces** with unexported marker methods:
- `chat.Message` (marker: `isMessage()`) — `SystemMessage`, `UserMessage`, `AssistantMessage`, `ToolMessage`
- `chat.ContentChunk` (marker: `contentChunk()`) — `TextChunk`, `ImageURLChunk`, `DocumentURLChunk`, `FileChunk`, `ReferenceChunk`, `ThinkChunk`, `AudioChunk`, `ToolReferenceChunk`, `ToolFileChunk`
- `agents.AgentTool` (marker: `agentToolType()`) — `FunctionTool`, `WebSearchTool`, `CodeInterpreterTool`, `ConnectorTool`, etc.
- `conversation.Event` — conversation streaming events

Each has an `Unknown*` variant so the SDK doesn't break on new API types. Each has a `Unmarshal*` dispatch function that probes a `type`/`role` discriminator field.

### Custom JSON patterns

Several types require non-trivial marshal/unmarshal:
- **Type alias trick** — `type alias T` inside `MarshalJSON` to avoid infinite recursion when injecting a `type`/`role` discriminator field.
- **`json:"-"` + custom MarshalJSON** — `CompletionRequest.Messages` (and `stream`) are excluded from default marshaling and injected via custom `MarshalJSON`.
- **Union types** — `Content` handles `string | null | []ContentChunk`; `ToolChoice` handles `string | object`; `ImageURL` handles `string | object`; `FunctionCall.Arguments` handles `string | object`; `ReferenceID` handles `int | string` with type preservation.
- **Probe struct pattern** — `Unmarshal*` functions decode only the discriminator field first, then dispatch to the concrete type.

### Shared types in `chat/`

`GuardrailConfig`, `ModerationLLMV1Config`, `ModerationLLMV2Config` live in `chat/` because it's the base types package imported by both `agents/` and `conversation/`. This avoids import cycles.

### Error handling

`APIError` in `error.go` with sentinel checkers: `IsNotFound()`, `IsRateLimit()`, `IsAuth()`. All use `errors.As` for unwrapping.

## Testing patterns

- Unit tests use `httptest.NewServer` with inline handlers to mock the Mistral API. Client is pointed at the test server via `WithBaseURL(server.URL)`.
- Integration tests are behind `//go:build integration` build tag and require `MISTRAL_API_KEY`.
- Tests use stdlib `testing` only — no third-party test frameworks.

## Adding a new API endpoint

1. Define request/response types in the appropriate sub-package (or create a new one with a `doc.go`).
2. Add a method on `*Client` in the root package. Use `doJSON` for standard request/response, `doStream` for SSE, `doMultipart` for file uploads.
3. Add unit tests with `httptest.NewServer`.
4. If the endpoint supports streaming, return `*Stream[T]` and call `EnableStream()` on the request before sending.
