package ocr

import "encoding/json"

// Request represents an OCR request.
type Request struct {
	Model                     *string         `json:"model"`
	ID                        string          `json:"id,omitempty"`
	Document                  json.RawMessage `json:"document"`
	Pages                     []int           `json:"pages,omitempty"`
	IncludeImageBase64        *bool           `json:"include_image_base64,omitempty"`
	ImageLimit                *int            `json:"image_limit,omitempty"`
	ImageMinSize              *int            `json:"image_min_size,omitempty"`
	BboxAnnotationFormat      json.RawMessage `json:"bbox_annotation_format,omitempty"`
	DocumentAnnotationFormat  json.RawMessage `json:"document_annotation_format,omitempty"`
	DocumentAnnotationPrompt  *string         `json:"document_annotation_prompt,omitempty"`
	TableFormat               *string         `json:"table_format,omitempty"`
	ExtractHeader             bool            `json:"extract_header,omitempty"`
	ExtractFooter             bool            `json:"extract_footer,omitempty"`
}

// Response is the OCR result.
type Response struct {
	Pages              []Page    `json:"pages"`
	Model              string    `json:"model"`
	DocumentAnnotation *string   `json:"document_annotation,omitempty"`
	UsageInfo          UsageInfo `json:"usage_info"`
}

// Page represents a single page's OCR results.
type Page struct {
	Index      int             `json:"index"`
	Markdown   string          `json:"markdown"`
	Images     []Image         `json:"images"`
	Tables     []Table         `json:"tables,omitempty"`
	Hyperlinks []string        `json:"hyperlinks,omitempty"`
	Header     *string         `json:"header,omitempty"`
	Footer     *string         `json:"footer,omitempty"`
	Dimensions *PageDimensions `json:"dimensions"`
}

// Image represents an extracted image from a page.
type Image struct {
	ID              string  `json:"id"`
	TopLeftX        *int    `json:"top_left_x"`
	TopLeftY        *int    `json:"top_left_y"`
	BottomRightX    *int    `json:"bottom_right_x"`
	BottomRightY    *int    `json:"bottom_right_y"`
	ImageBase64     *string `json:"image_base64,omitempty"`
	ImageAnnotation *string `json:"image_annotation,omitempty"`
}

// Table represents an extracted table from a page.
type Table struct {
	ID      string `json:"id"`
	Content string `json:"content"`
	Format  string `json:"format"`
}

// PageDimensions holds the dimensions of a page image.
type PageDimensions struct {
	DPI    int `json:"dpi"`
	Height int `json:"height"`
	Width  int `json:"width"`
}

// UsageInfo holds OCR usage statistics.
type UsageInfo struct {
	PagesProcessed int  `json:"pages_processed"`
	DocSizeBytes   *int `json:"doc_size_bytes,omitempty"`
}
