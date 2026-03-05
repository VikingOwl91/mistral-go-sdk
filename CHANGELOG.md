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
