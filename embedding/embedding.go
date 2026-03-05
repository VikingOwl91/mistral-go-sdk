package embedding

import "somegit.dev/vikingowl/mistral-go-sdk/chat"

// Dtype specifies the data type of output embeddings.
type Dtype string

const (
	DtypeFloat   Dtype = "float"
	DtypeInt8    Dtype = "int8"
	DtypeUint8   Dtype = "uint8"
	DtypeBinary  Dtype = "binary"
	DtypeUbinary Dtype = "ubinary"
)

// EncodingFormat specifies the format of embeddings in the response.
type EncodingFormat string

const (
	EncodingFormatFloat  EncodingFormat = "float"
	EncodingFormatBase64 EncodingFormat = "base64"
)

// Request represents an embedding request.
type Request struct {
	Model           string          `json:"model"`
	Input           []string        `json:"input"`
	OutputDimension *int            `json:"output_dimension,omitempty"`
	OutputDtype     *Dtype          `json:"output_dtype,omitempty"`
	EncodingFormat  *EncodingFormat `json:"encoding_format,omitempty"`
	Metadata        map[string]any  `json:"metadata,omitempty"`
}

// Response represents an embedding response.
type Response struct {
	ID    string         `json:"id"`
	Object string        `json:"object"`
	Model  string        `json:"model"`
	Usage  chat.UsageInfo `json:"usage"`
	Data   []Data         `json:"data"`
}

// Data represents a single embedding result.
type Data struct {
	Object    string    `json:"object"`
	Embedding []float64 `json:"embedding"`
	Index     int       `json:"index"`
}
