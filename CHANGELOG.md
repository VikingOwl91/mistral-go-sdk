## v1.4.0 — 2026-04-28

Spec/SDK alignment pass after upstream OpenAPI moved to v1.0.0 and
Python SDK shipped v2.3.0..v2.4.3. RAG ingestion-pipeline beta surface
(Python v2.4.3) intentionally deferred until the dust settles upstream.

### Added

- **`Client.GetWorkflowWorkerInfo`** — restores the
  `GET /v1/workflows/workers/whoami` endpoint that was removed in v1.3.0.
  The endpoint is still in the spec and is needed by callers running
  custom workers that connect their own scheduler.
  (`workflow.WorkerInfo` type.)
- **Observability fields API** — three GETs missing since the
  observability surface was first added:
  - `Client.GetChatCompletionFields` (`/v1/observability/chat-completion-fields`)
  - `Client.GetChatCompletionFieldOptions` (`…/{field}/options?operator=…`)
  - `Client.GetChatCompletionFieldOptionsCounts` (`…/{field}/options-counts`)
  - new types: `observability.BaseFieldDefinition`, `FieldGroup`,
    `ChatCompletionFields`, `ChatCompletionFieldOptions`,
    `FieldOptionCountsRequest`, `FieldOptionCounts`, `FieldOptionCountItem`,
    plus `FieldType` and `FieldOperator` typed enums.
- **Workflow payload encoding constants** — `workflow.EncodedPayloadOption`
  with `EncodedPayloadOffloaded`, `EncodedPayloadEncrypted`,
  `EncodedPayloadEncryptedPartial`. Wire-compatible refinement of the
  pre-existing `[]string` field on `NetworkEncodedInput`.
  (Mirrors Python SDK v2.4.0.)
- **Workflow ↔ connector integration** (Python SDK v2.4.2):
  - `workflow.ConnectorSlot`, `ConnectorBindings`, `ConnectorExtensions`,
    `WorkflowExtensions` types.
  - `workflow.BuildConnectorExtensions(slots …)` helper that produces the
    nested map expected at `ExecutionRequest.Extensions["mistralai"]`.
  - `workflow.ConnectorAuthTaskState` + `ConnectorAuthStatus` constants
    for parsing payloads emitted by the `connector-auth` custom task event.
  - New `Extensions map[string]any` field on `workflow.ExecutionRequest`.
- **HITL (human-in-the-loop) confirmation constants** — typed values
  alongside the pre-existing `conversation.ToolCallConfirmation` and
  `tool_confirmations` field:
  - `conversation.Confirmation` with `ConfirmationAllow` / `ConfirmationDeny`
    for the reply side.
  - `ConfirmationStatusPending` / `ConfirmationStatusAllowed` /
    `ConfirmationStatusDenied` for `FunctionCallEvent.ConfirmationStatus`
    and `FunctionCallEntry.ConfirmationStatus` (already present as
    untyped strings).

### Changed

- `workflow.NetworkEncodedInput.EncodingOptions` is now
  `[]EncodedPayloadOption` (string-typed alias). JSON wire format
  unchanged; existing call sites that passed `[]string{"offloaded"}`
  need to switch to `[]workflow.EncodedPayloadOption{workflow.EncodedPayloadOffloaded}`
  or the typed constants directly.
- Tracking upstream Mistral OpenAPI spec **v1.0.0** (was v0.1.104).
  Only spec delta in this window was the removal of OCR confidence-score
  fields (`OCRPageObject.confidence_scores`,
  `OCRRequest.confidence_scores_granularity`,
  `OCRTableObject.word_confidence_scores`, plus the `OCRConfidenceScore`
  and `OCRPageConfidenceScores` schemas), none of which this SDK exposed.

### Fixed (CI)

- `watch-openapi.yml` no longer attempts to commit `.openapi-hash` /
  `.openapi-spec.yaml` to `main`. The push was being silently reverted
  by an upstream mirror, leaving the tracking issue stale across
  multiple upstream releases. The watcher now refreshes the open
  tracking issue's body on each run so the diff and hashes always
  reflect the current upstream state, and posts a comment when the
  spec moves again while the issue is still open.

## v1.3.0 — 2026-04-03

Upstream sync with Python SDK v2.3.0. Updates workflow registration model
to reflect the managed deployment architecture and removes the deprecated
workers endpoint.

### Added

- **`workflow.CodeDefinition`** — workflow interface metadata type with
  input/output schemas, signal/query/update handler definitions,
  determinism flag, and execution timeout.
- **`workflow.SignalDefinition`**, **`QueryDefinition`**,
  **`UpdateDefinition`** — handler descriptor types.
- **`Registration.Definition`** — code definition field on workflow
  registrations.
- **`Registration.DeploymentID`** — replaces the worker/task-queue model
  with managed deployment references.
- **`Registration.CompatibleWithChatAssistant`** — flag for chat assistant
  compatibility.

### Deprecated

- **`Registration.TaskQueue`** — use `DeploymentID` instead. Will be
  removed in a future release.

### Removed (breaking)

- **`GetWorkflowWorkerInfo`** — the `/v1/workflows/workers/whoami` endpoint
  was removed upstream.
- **`workflow.WorkerInfo`** — type no longer exists in the API.

## v1.2.1 — 2026-04-03

Move module path to `github.com/VikingOwl91/mistral-go-sdk` for public
discoverability on pkg.go.dev.

### Changed

- Module path changed from `somegit.dev/vikingowl/mistral-go-sdk` to
  `github.com/VikingOwl91/mistral-go-sdk`.

### Fixed

- `TestChatCompleteStream_WithToolCalls` fixture now includes `finish_reason`
  and `usage` to match real Mistral API responses.

## v1.2.0 — 2026-04-02

Upstream sync with Python SDK v2.2.0. Adds Workflows API and DeleteBatchJob.

### Added

- **Workflows API** (new `workflow/` package) — complete workflow orchestration
  support with 37 service methods across 8 sub-resources:
  - **Workflows CRUD** — `ListWorkflows`, `GetWorkflow`, `UpdateWorkflow`,
    `ArchiveWorkflow`, `UnarchiveWorkflow`, `ExecuteWorkflow`,
    `ExecuteWorkflowAndWait`.
  - **Registrations** — `ListWorkflowRegistrations`, `GetWorkflowRegistration`,
    `ExecuteWorkflowRegistration` (deprecated).
  - **Executions** — `GetWorkflowExecution`, `GetWorkflowExecutionHistory`,
    `StreamWorkflowExecution`, `SignalWorkflowExecution`,
    `QueryWorkflowExecution`, `UpdateWorkflowExecution`,
    `TerminateWorkflowExecution`, `CancelWorkflowExecution`,
    `ResetWorkflowExecution`, `BatchCancelWorkflowExecutions`,
    `BatchTerminateWorkflowExecutions`.
  - **Trace** — `GetWorkflowExecutionTraceOTel`,
    `GetWorkflowExecutionTraceSummary`, `GetWorkflowExecutionTraceEvents`.
  - **Events** — `StreamWorkflowEvents`, `ListWorkflowEvents`.
  - **Deployments** — `ListWorkflowDeployments`, `GetWorkflowDeployment`.
  - **Metrics** — `GetWorkflowMetrics`.
  - **Runs** — `ListWorkflowRuns`, `GetWorkflowRun`, `GetWorkflowRunHistory`.
  - **Schedules** — `ListWorkflowSchedules`, `ScheduleWorkflow`,
    `UnscheduleWorkflow`.
  - **Workers** — `GetWorkflowWorkerInfo`.
- **`WorkflowEventStream`** — typed SSE stream wrapper with `StreamPayload`
  envelope, sealed `Event` interface (17 concrete types + `UnknownEvent`).
- **`DeleteBatchJob`** — delete a batch job by ID.

## v1.1.0 — 2026-03-24

Upstream sync with Python SDK v2.1.3. Adds Connectors, Audio Speech/Voices, and Observability (beta).

### Breaking Changes

- **`ListModels`** signature changed from `(ctx)` to `(ctx, *model.ListParams)`.
  Pass `nil` for previous behavior. The new `ListParams` supports `Provider` and
  `Model` query filters.
- **`UploadFile`** signature changed from `(ctx, filename, reader, purpose)` to
  `(ctx, filename, reader, *file.UploadParams)`. The new `UploadParams` struct
  holds `Purpose`, `Expiry`, and `Visibility` fields.

### Added

- **`ReasoningEffort`** field on `chat.CompletionRequest` and
  `agents.CompletionRequest` — controls reasoning effort (`"none"`, `"high"`).
- **Connectors API** (new `connector/` package) — `CreateConnector`,
  `ListConnectors`, `GetConnector`, `UpdateConnector`, `DeleteConnector`,
  `GetConnectorAuthURL`, `ListConnectorTools`, `CallConnectorTool`.
- **Audio Speech (TTS)** — `Speech`, `SpeechStream` with `SpeechStream` typed
  wrapper, `SpeechOutputFormat` enum (pcm/wav/mp3/flac/opus).
- **Audio Voices** — `ListVoices`, `CreateVoice`, `GetVoice`, `UpdateVoice`,
  `DeleteVoice`, `GetVoiceSampleAudio`.
- **Audio Realtime types** — `AudioEncoding`, `AudioFormat`, `RealtimeSession`,
  and WebSocket message types in `audio/realtime.go`. No WebSocket client yet
  (would require adding a dependency).
- **Observability API** (new `observability/` package, beta) — campaigns,
  chat completion events, judges, datasets, records, and import tasks.
  33 service methods total.
- **`file.Visibility`** enum — `shared_global`, `shared_org`,
  `shared_workspace`, `private`.
- **`model.ListParams`** — filter models by `Provider` and `Model`.

## v1.0.0 — 2026-03-17

Stable release. Tracks upstream Python SDK v2.0.4.

No API changes from v0.2.0. This release signals that the SDK surface is
stable and follows Go module semver conventions — breaking changes will only
occur in future major versions.

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
