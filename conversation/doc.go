// Package conversation provides types for the Mistral conversations API.
//
// Conversations support multi-turn interactions with start, append, and
// restart operations. Streaming returns typed [Event] values via a sealed
// interface dispatched by the "type" discriminator.
package conversation
