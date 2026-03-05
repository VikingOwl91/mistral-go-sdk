package chat

import (
	"encoding/json"
	"testing"
)

func TestTextChunk_MarshalJSON(t *testing.T) {
	chunk := &TextChunk{Text: "hello"}
	data, err := json.Marshal(chunk)
	if err != nil {
		t.Fatal(err)
	}
	want := `{"type":"text","text":"hello"}`
	if string(data) != want {
		t.Errorf("got %s, want %s", data, want)
	}
}

func TestTextChunk_RoundTrip(t *testing.T) {
	original := &TextChunk{Text: "hello world"}
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatal(err)
	}
	chunk, err := UnmarshalContentChunk(data)
	if err != nil {
		t.Fatal(err)
	}
	tc, ok := chunk.(*TextChunk)
	if !ok {
		t.Fatalf("expected *TextChunk, got %T", chunk)
	}
	if tc.Text != original.Text {
		t.Errorf("got %q, want %q", tc.Text, original.Text)
	}
}

func TestImageURLChunk_MarshalJSON(t *testing.T) {
	chunk := &ImageURLChunk{ImageURL: ImageURL{URL: "https://example.com/img.png"}}
	data, err := json.Marshal(chunk)
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]any
	json.Unmarshal(data, &m)
	if m["type"] != "image_url" {
		t.Errorf("expected type=image_url, got %v", m["type"])
	}
	imgURL, ok := m["image_url"].(map[string]any)
	if !ok {
		t.Fatalf("expected image_url object, got %T", m["image_url"])
	}
	if imgURL["url"] != "https://example.com/img.png" {
		t.Errorf("expected url, got %v", imgURL["url"])
	}
}

func TestImageURLChunk_RoundTrip(t *testing.T) {
	detail := ImageDetailHigh
	original := &ImageURLChunk{ImageURL: ImageURL{
		URL:    "https://example.com/img.png",
		Detail: &detail,
	}}
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatal(err)
	}
	chunk, err := UnmarshalContentChunk(data)
	if err != nil {
		t.Fatal(err)
	}
	ic, ok := chunk.(*ImageURLChunk)
	if !ok {
		t.Fatalf("expected *ImageURLChunk, got %T", chunk)
	}
	if ic.ImageURL.URL != original.ImageURL.URL {
		t.Errorf("got URL %q, want %q", ic.ImageURL.URL, original.ImageURL.URL)
	}
	if ic.ImageURL.Detail == nil || *ic.ImageURL.Detail != ImageDetailHigh {
		t.Errorf("got detail %v, want high", ic.ImageURL.Detail)
	}
}

func TestImageURL_UnmarshalJSON_String(t *testing.T) {
	var u ImageURL
	if err := json.Unmarshal([]byte(`"https://example.com/img.png"`), &u); err != nil {
		t.Fatal(err)
	}
	if u.URL != "https://example.com/img.png" {
		t.Errorf("got %q, want %q", u.URL, "https://example.com/img.png")
	}
}

func TestImageURL_UnmarshalJSON_Object(t *testing.T) {
	var u ImageURL
	if err := json.Unmarshal([]byte(`{"url":"https://example.com/img.png","detail":"low"}`), &u); err != nil {
		t.Fatal(err)
	}
	if u.URL != "https://example.com/img.png" {
		t.Errorf("got url %q", u.URL)
	}
	if u.Detail == nil || *u.Detail != ImageDetailLow {
		t.Errorf("got detail %v, want low", u.Detail)
	}
}

func TestDocumentURLChunk_RoundTrip(t *testing.T) {
	name := "doc.pdf"
	original := &DocumentURLChunk{
		DocumentURL:  "https://example.com/doc.pdf",
		DocumentName: &name,
	}
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatal(err)
	}
	chunk, err := UnmarshalContentChunk(data)
	if err != nil {
		t.Fatal(err)
	}
	dc, ok := chunk.(*DocumentURLChunk)
	if !ok {
		t.Fatalf("expected *DocumentURLChunk, got %T", chunk)
	}
	if dc.DocumentURL != original.DocumentURL {
		t.Errorf("got %q, want %q", dc.DocumentURL, original.DocumentURL)
	}
	if dc.DocumentName == nil || *dc.DocumentName != name {
		t.Errorf("got name %v, want %q", dc.DocumentName, name)
	}
}

func TestFileChunk_RoundTrip(t *testing.T) {
	original := &FileChunk{FileID: "abc-123"}
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatal(err)
	}
	chunk, err := UnmarshalContentChunk(data)
	if err != nil {
		t.Fatal(err)
	}
	fc, ok := chunk.(*FileChunk)
	if !ok {
		t.Fatalf("expected *FileChunk, got %T", chunk)
	}
	if fc.FileID != original.FileID {
		t.Errorf("got %q, want %q", fc.FileID, original.FileID)
	}
}

func TestReferenceChunk_RoundTrip(t *testing.T) {
	original := &ReferenceChunk{ReferenceIDs: []int{1, 2, 3}}
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatal(err)
	}
	chunk, err := UnmarshalContentChunk(data)
	if err != nil {
		t.Fatal(err)
	}
	rc, ok := chunk.(*ReferenceChunk)
	if !ok {
		t.Fatalf("expected *ReferenceChunk, got %T", chunk)
	}
	if len(rc.ReferenceIDs) != 3 || rc.ReferenceIDs[0] != 1 {
		t.Errorf("got %v, want [1 2 3]", rc.ReferenceIDs)
	}
}

func TestThinkChunk_RoundTrip(t *testing.T) {
	closed := true
	original := &ThinkChunk{
		Thinking: []ContentChunk{
			&TextChunk{Text: "reasoning step"},
			&ReferenceChunk{ReferenceIDs: []int{42}},
		},
		Closed: &closed,
	}
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatal(err)
	}
	chunk, err := UnmarshalContentChunk(data)
	if err != nil {
		t.Fatal(err)
	}
	tc, ok := chunk.(*ThinkChunk)
	if !ok {
		t.Fatalf("expected *ThinkChunk, got %T", chunk)
	}
	if len(tc.Thinking) != 2 {
		t.Fatalf("got %d thinking chunks, want 2", len(tc.Thinking))
	}
	if tc.Closed == nil || !*tc.Closed {
		t.Error("expected closed=true")
	}
	text, ok := tc.Thinking[0].(*TextChunk)
	if !ok {
		t.Fatalf("expected thinking[0] to be *TextChunk, got %T", tc.Thinking[0])
	}
	if text.Text != "reasoning step" {
		t.Errorf("got %q, want %q", text.Text, "reasoning step")
	}
}

func TestAudioChunk_RoundTrip(t *testing.T) {
	original := &AudioChunk{InputAudio: "base64data=="}
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatal(err)
	}
	chunk, err := UnmarshalContentChunk(data)
	if err != nil {
		t.Fatal(err)
	}
	ac, ok := chunk.(*AudioChunk)
	if !ok {
		t.Fatalf("expected *AudioChunk, got %T", chunk)
	}
	if ac.InputAudio != original.InputAudio {
		t.Errorf("got %q, want %q", ac.InputAudio, original.InputAudio)
	}
}

func TestUnmarshalContentChunk_UnknownType(t *testing.T) {
	data := []byte(`{"type":"future_chunk","data":"value"}`)
	chunk, err := UnmarshalContentChunk(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	u, ok := chunk.(*UnknownChunk)
	if !ok {
		t.Fatalf("expected *UnknownChunk, got %T", chunk)
	}
	if u.Type != "future_chunk" {
		t.Errorf("got type %q", u.Type)
	}
	marshaled, err := json.Marshal(u)
	if err != nil {
		t.Fatal(err)
	}
	if string(marshaled) != string(data) {
		t.Errorf("round-trip failed: got %s", marshaled)
	}
}

func TestContent_TextMarshal(t *testing.T) {
	c := TextContent("hello")
	data, err := json.Marshal(c)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != `"hello"` {
		t.Errorf("got %s, want %q", data, `"hello"`)
	}
}

func TestContent_ChunksMarshal(t *testing.T) {
	c := ChunksContent(&TextChunk{Text: "hello"})
	data, err := json.Marshal(c)
	if err != nil {
		t.Fatal(err)
	}
	want := `[{"type":"text","text":"hello"}]`
	if string(data) != want {
		t.Errorf("got %s, want %s", data, want)
	}
}

func TestContent_NullMarshal(t *testing.T) {
	var c Content
	data, err := json.Marshal(c)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "null" {
		t.Errorf("got %s, want null", data)
	}
}

func TestContent_UnmarshalString(t *testing.T) {
	var c Content
	if err := json.Unmarshal([]byte(`"hello"`), &c); err != nil {
		t.Fatal(err)
	}
	if c.Text == nil || *c.Text != "hello" {
		t.Errorf("expected text 'hello', got %+v", c)
	}
	if c.String() != "hello" {
		t.Errorf("String() = %q, want %q", c.String(), "hello")
	}
}

func TestContent_UnmarshalArray(t *testing.T) {
	var c Content
	if err := json.Unmarshal([]byte(`[{"type":"text","text":"hi"}]`), &c); err != nil {
		t.Fatal(err)
	}
	if len(c.Parts) != 1 {
		t.Fatalf("got %d parts, want 1", len(c.Parts))
	}
	tc, ok := c.Parts[0].(*TextChunk)
	if !ok {
		t.Fatalf("expected *TextChunk, got %T", c.Parts[0])
	}
	if tc.Text != "hi" {
		t.Errorf("got %q, want %q", tc.Text, "hi")
	}
}

func TestContent_UnmarshalNull(t *testing.T) {
	var c Content
	if err := json.Unmarshal([]byte("null"), &c); err != nil {
		t.Fatal(err)
	}
	if !c.IsNull() {
		t.Error("expected null content")
	}
}

func TestContent_RoundTrip_Text(t *testing.T) {
	original := TextContent("round trip")
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatal(err)
	}
	var decoded Content
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded.String() != "round trip" {
		t.Errorf("got %q, want %q", decoded.String(), "round trip")
	}
}

func TestContent_RoundTrip_Chunks(t *testing.T) {
	original := ChunksContent(
		&TextChunk{Text: "hello"},
		&ImageURLChunk{ImageURL: ImageURL{URL: "https://example.com/img.png"}},
	)
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatal(err)
	}
	var decoded Content
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatal(err)
	}
	if len(decoded.Parts) != 2 {
		t.Fatalf("got %d parts, want 2", len(decoded.Parts))
	}
}

func TestContent_IsNull(t *testing.T) {
	var zero Content
	if !zero.IsNull() {
		t.Error("zero value should be null")
	}
	if TextContent("x").IsNull() {
		t.Error("text content should not be null")
	}
	if ChunksContent(&TextChunk{Text: "x"}).IsNull() {
		t.Error("chunks content should not be null")
	}
}
