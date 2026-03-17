package chat

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// BuiltInConnector identifies a built-in connector type.
type BuiltInConnector string

const (
	ConnectorWebSearch        BuiltInConnector = "web_search"
	ConnectorWebSearchPremium BuiltInConnector = "web_search_premium"
	ConnectorCodeInterpreter  BuiltInConnector = "code_interpreter"
	ConnectorImageGeneration  BuiltInConnector = "image_generation"
	ConnectorDocumentLibrary  BuiltInConnector = "document_library"
)

// ContentChunk is a sealed interface for message content parts.
// Only concrete types in this package can implement it.
type ContentChunk interface {
	contentChunk()
}

// TextChunk represents a text content part.
type TextChunk struct {
	Text string `json:"text"`
}

func (*TextChunk) contentChunk() {}

func (c *TextChunk) MarshalJSON() ([]byte, error) {
	type alias TextChunk
	return json.Marshal(&struct {
		Type string `json:"type"`
		*alias
	}{
		Type:  "text",
		alias: (*alias)(c),
	})
}

// ImageDetail controls image processing fidelity.
type ImageDetail string

const (
	ImageDetailLow  ImageDetail = "low"
	ImageDetailAuto ImageDetail = "auto"
	ImageDetailHigh ImageDetail = "high"
)

// ImageURL holds the URL and optional detail level for an image.
type ImageURL struct {
	URL    string       `json:"url"`
	Detail *ImageDetail `json:"detail,omitempty"`
}

func (u *ImageURL) UnmarshalJSON(data []byte) error {
	if data[0] == '"' {
		var s string
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}
		u.URL = s
		return nil
	}
	type alias ImageURL
	return json.Unmarshal(data, (*alias)(u))
}

// ImageURLChunk represents an image content part.
type ImageURLChunk struct {
	ImageURL ImageURL `json:"image_url"`
}

func (*ImageURLChunk) contentChunk() {}

func (c *ImageURLChunk) MarshalJSON() ([]byte, error) {
	type alias ImageURLChunk
	return json.Marshal(&struct {
		Type string `json:"type"`
		*alias
	}{
		Type:  "image_url",
		alias: (*alias)(c),
	})
}

// DocumentURLChunk represents a document URL content part.
type DocumentURLChunk struct {
	DocumentURL  string  `json:"document_url"`
	DocumentName *string `json:"document_name,omitempty"`
}

func (*DocumentURLChunk) contentChunk() {}

func (c *DocumentURLChunk) MarshalJSON() ([]byte, error) {
	type alias DocumentURLChunk
	return json.Marshal(&struct {
		Type string `json:"type"`
		*alias
	}{
		Type:  "document_url",
		alias: (*alias)(c),
	})
}

// FileChunk represents a file reference content part.
type FileChunk struct {
	FileID string `json:"file_id"`
}

func (*FileChunk) contentChunk() {}

func (c *FileChunk) MarshalJSON() ([]byte, error) {
	type alias FileChunk
	return json.Marshal(&struct {
		Type string `json:"type"`
		*alias
	}{
		Type:  "file",
		alias: (*alias)(c),
	})
}

// ReferenceID is a reference identifier that can be an integer or string.
// Use [IntRef] or [StringRef] constructors.
type ReferenceID struct {
	raw      string
	isString bool
}

// IntRef creates an integer reference ID.
func IntRef(n int) ReferenceID {
	return ReferenceID{raw: strconv.Itoa(n)}
}

// StringRef creates a string reference ID.
func StringRef(s string) ReferenceID {
	return ReferenceID{raw: s, isString: true}
}

// String returns the string representation.
func (id ReferenceID) String() string { return id.raw }

// Int returns the integer value and true if this is a numeric reference.
func (id ReferenceID) Int() (int, bool) {
	if id.isString {
		return 0, false
	}
	n, err := strconv.Atoi(id.raw)
	return n, err == nil
}

// IsString reports whether this is a string reference.
func (id ReferenceID) IsString() bool { return id.isString }

func (id ReferenceID) MarshalJSON() ([]byte, error) {
	if id.isString {
		return json.Marshal(id.raw)
	}
	return []byte(id.raw), nil
}

func (id *ReferenceID) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		return nil
	}
	if data[0] == '"' {
		var s string
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}
		id.raw = s
		id.isString = true
		return nil
	}
	id.raw = string(data)
	return nil
}

// ReferenceChunk represents a reference content part.
type ReferenceChunk struct {
	ReferenceIDs []ReferenceID `json:"reference_ids"`
}

func (*ReferenceChunk) contentChunk() {}

func (c *ReferenceChunk) MarshalJSON() ([]byte, error) {
	type alias ReferenceChunk
	return json.Marshal(&struct {
		Type string `json:"type"`
		*alias
	}{
		Type:  "reference",
		alias: (*alias)(c),
	})
}

// ThinkChunk represents a thinking/reasoning content part.
type ThinkChunk struct {
	Thinking []ContentChunk `json:"-"`
	Closed   *bool          `json:"closed,omitempty"`
}

func (*ThinkChunk) contentChunk() {}

func (c *ThinkChunk) MarshalJSON() ([]byte, error) {
	thinking, err := json.Marshal(c.Thinking)
	if err != nil {
		return nil, err
	}
	return json.Marshal(&struct {
		Type     string          `json:"type"`
		Thinking json.RawMessage `json:"thinking"`
		Closed   *bool           `json:"closed,omitempty"`
	}{
		Type:     "thinking",
		Thinking: thinking,
		Closed:   c.Closed,
	})
}

func (c *ThinkChunk) UnmarshalJSON(data []byte) error {
	var raw struct {
		Thinking []json.RawMessage `json:"thinking"`
		Closed   *bool             `json:"closed"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	c.Closed = raw.Closed
	c.Thinking = make([]ContentChunk, len(raw.Thinking))
	for i, r := range raw.Thinking {
		chunk, err := UnmarshalContentChunk(r)
		if err != nil {
			return err
		}
		c.Thinking[i] = chunk
	}
	return nil
}

// AudioChunk represents an audio input content part.
type AudioChunk struct {
	InputAudio string `json:"input_audio"`
}

func (*AudioChunk) contentChunk() {}

func (c *AudioChunk) MarshalJSON() ([]byte, error) {
	type alias AudioChunk
	return json.Marshal(&struct {
		Type string `json:"type"`
		*alias
	}{
		Type:  "input_audio",
		alias: (*alias)(c),
	})
}

// ToolReferenceChunk represents a tool reference content part.
type ToolReferenceChunk struct {
	Tool        string  `json:"tool"`
	Title       string  `json:"title"`
	URL         *string `json:"url,omitempty"`
	Favicon     *string `json:"favicon,omitempty"`
	Description *string `json:"description,omitempty"`
}

func (*ToolReferenceChunk) contentChunk() {}

func (c *ToolReferenceChunk) MarshalJSON() ([]byte, error) {
	type alias ToolReferenceChunk
	return json.Marshal(&struct {
		Type string `json:"type"`
		*alias
	}{
		Type:  "tool_reference",
		alias: (*alias)(c),
	})
}

// ToolFileChunk represents a tool-generated file content part.
type ToolFileChunk struct {
	Tool     string  `json:"tool"`
	FileID   string  `json:"file_id"`
	FileName *string `json:"file_name,omitempty"`
	FileType *string `json:"file_type,omitempty"`
}

func (*ToolFileChunk) contentChunk() {}

func (c *ToolFileChunk) MarshalJSON() ([]byte, error) {
	type alias ToolFileChunk
	return json.Marshal(&struct {
		Type string `json:"type"`
		*alias
	}{
		Type:  "tool_file",
		alias: (*alias)(c),
	})
}

// UnmarshalContentChunk dispatches to the concrete ContentChunk type
// based on the "type" discriminator field.
func UnmarshalContentChunk(data []byte) (ContentChunk, error) {
	var probe struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(data, &probe); err != nil {
		return nil, fmt.Errorf("mistral: unmarshal content chunk: %w", err)
	}
	switch probe.Type {
	case "text":
		var c TextChunk
		return &c, json.Unmarshal(data, &c)
	case "image_url":
		var c ImageURLChunk
		return &c, json.Unmarshal(data, &c)
	case "document_url":
		var c DocumentURLChunk
		return &c, json.Unmarshal(data, &c)
	case "file":
		var c FileChunk
		return &c, json.Unmarshal(data, &c)
	case "reference":
		var c ReferenceChunk
		return &c, json.Unmarshal(data, &c)
	case "thinking":
		var c ThinkChunk
		return &c, json.Unmarshal(data, &c)
	case "input_audio":
		var c AudioChunk
		return &c, json.Unmarshal(data, &c)
	case "tool_reference":
		var c ToolReferenceChunk
		return &c, json.Unmarshal(data, &c)
	case "tool_file":
		var c ToolFileChunk
		return &c, json.Unmarshal(data, &c)
	default:
		return &UnknownChunk{Type: probe.Type, Raw: json.RawMessage(data)}, nil
	}
}

// UnknownChunk holds a content chunk with an unrecognized type.
// This prevents the SDK from breaking when new chunk types are added.
type UnknownChunk struct {
	Type string
	Raw  json.RawMessage
}

func (*UnknownChunk) contentChunk() {}

func (c *UnknownChunk) MarshalJSON() ([]byte, error) {
	return c.Raw, nil
}

// Content represents a message content field that can be a string,
// an array of content chunks, or null.
type Content struct {
	Text  *string
	Parts []ContentChunk
}

// TextContent creates a Content holding a plain string.
func TextContent(s string) Content {
	return Content{Text: &s}
}

// ChunksContent creates a Content holding an array of content chunks.
func ChunksContent(chunks ...ContentChunk) Content {
	return Content{Parts: chunks}
}

// IsNull returns true if the content is null (neither text nor chunks).
func (c Content) IsNull() bool {
	return c.Text == nil && c.Parts == nil
}

// String returns the text content, or empty string if not text.
func (c Content) String() string {
	if c.Text != nil {
		return *c.Text
	}
	return ""
}

func (c Content) MarshalJSON() ([]byte, error) {
	if c.Text != nil {
		return json.Marshal(*c.Text)
	}
	if c.Parts != nil {
		return json.Marshal(c.Parts)
	}
	return []byte("null"), nil
}

func (c *Content) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	if len(data) == 0 {
		return nil
	}
	if data[0] == '"' {
		var s string
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}
		c.Text = &s
		return nil
	}
	if data[0] == '[' {
		var raw []json.RawMessage
		if err := json.Unmarshal(data, &raw); err != nil {
			return err
		}
		c.Parts = make([]ContentChunk, len(raw))
		for i, r := range raw {
			chunk, err := UnmarshalContentChunk(r)
			if err != nil {
				return err
			}
			c.Parts[i] = chunk
		}
		return nil
	}
	return fmt.Errorf("mistral: content must be a string, array, or null")
}
