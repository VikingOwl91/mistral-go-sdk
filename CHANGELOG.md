## v0.2.0 — 2026-03-17

Sync with upstream Python SDK v2.0.4. Upstream reference changed from OpenAPI
spec to official Python SDK (https://github.com/mistralai/client-python).

### Breaking Changes

- **`ReferenceChunk.ReferenceIDs`** changed from `[]int` to `[]ReferenceID`.
  The API now returns mixed integer and string identifiers. Use `IntRef(n)` and
  `StringRef(s)` constructors; read back with `.Int()` and `.IsString()`.
- **`agents.GuardrailConfig` and `agents.ModerationLLMV1Config`** moved to
  `chat.GuardrailConfig` and `chat.ModerationLLMV1Config`. The types are now
  shared across chat, agents, and conversation packages.

### Added

- **`ToolReferenceChunk`** — new content chunk type for tool references
  returned by built-in connectors (web search, code interpreter, etc.).
- **`ToolFileChunk`** — new content chunk type for tool-generated files.
- **`BuiltInConnector`** constants — `ConnectorWebSearch`,
  `ConnectorWebSearchPremium`, `ConnectorCodeInterpreter`,
  `ConnectorImageGeneration`, `ConnectorDocumentLibrary`.
- **`ModerationLLMV2Config`** — v2 moderation guardrail with split
  `dangerous`/`criminal` categories and new `jailbreaking` category.
- **`GuardrailConfig`** on `chat.CompletionRequest`,
  `agents.CompletionRequest`, `conversation.StartRequest`, and
  `conversation.RestartRequest`.
- **`ConnectorTool`** — new agent tool type for custom connectors with
  `ConnectorAuth` (api-key / oauth2-token authorization).
- **`ModelCapabilities`** — added `AudioTranscriptionRealtime` and
  `AudioSpeech` fields.

### Removed

- Bundled `docs/openapi.yaml`. The SDK now tracks the upstream Python SDK
  directly as its reference implementation.

## v0.1.0 — 2026-03-05

Initial release.

- Chat completions (sync + streaming)
- Fill-in-the-Middle (FIM) completions
- Agents completions and CRUD (create, list, get, update, delete, versions, aliases)
- Conversations (start, append, restart, stream, history, messages)
- Embeddings
- Models (list, get, delete)
- Files (upload, list, get, delete, signed URLs)
- Fine-tuning jobs (create, list, get, cancel, start, archive/unarchive)
- Batch jobs (create, list, get, cancel)
- OCR (document processing)
- Audio transcription (sync + streaming)
- Document libraries (CRUD, upload, sharing, reprocessing)
- Moderation and classification
- Retry with exponential backoff, jitter, and Retry-After support
- Typed streaming via generic `Stream[T]` iterator
- Sealed interface pattern for discriminated unions (messages, events, entries)
