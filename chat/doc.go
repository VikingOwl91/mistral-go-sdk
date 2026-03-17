// Package chat provides types for the Mistral chat completion API.
//
// Messages use a sealed interface pattern — create them with
// [UserMessage], [SystemMessage], [AssistantMessage], or [ToolMessage]
// struct literals.
//
// Content is polymorphic: it can be a plain string (via [TextContent]),
// nil, or a slice of [ContentChunk] values (text, image URL, document URL,
// audio, tool reference, tool file).
//
// # Guardrails
//
// [GuardrailConfig] configures moderation guardrails on completion requests.
// Both v1 ([ModerationLLMV1Config]) and v2 ([ModerationLLMV2Config]) moderation
// configs are supported.
//
// # Content Chunks
//
// The following chunk types are supported: [TextChunk], [ImageURLChunk],
// [DocumentURLChunk], [FileChunk], [ReferenceChunk], [ThinkChunk],
// [AudioChunk], [ToolReferenceChunk], [ToolFileChunk].
// Unrecognized types are preserved as [UnknownChunk].
//
// [ReferenceChunk] uses [ReferenceID] values that can hold either integer
// or string identifiers. Use [IntRef] and [StringRef] constructors.
package chat
