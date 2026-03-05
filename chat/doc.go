// Package chat provides types for the Mistral chat completion API.
//
// Messages use a sealed interface pattern — create them with
// [UserMessage], [SystemMessage], [AssistantMessage], or [ToolMessage]
// struct literals.
//
// Content is polymorphic: it can be a plain string (via [TextContent]),
// nil, or a slice of [ContentChunk] values (text, image URL, document URL, audio).
package chat
