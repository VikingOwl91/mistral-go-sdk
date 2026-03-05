// Package model provides types for Mistral model listing and metadata.
//
// Model cards use a sealed interface — [BaseModelCard] for platform models
// and [FTModelCard] for fine-tuned models. JSON unmarshaling dispatches on
// the "type" field.
package model
