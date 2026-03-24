package file

// Purpose indicates the intended use of an uploaded file.
type Purpose string

const (
	PurposeFineTune Purpose = "fine-tune"
	PurposeBatch    Purpose = "batch"
	PurposeOCR      Purpose = "ocr"
)

// SampleType categorizes file content.
type SampleType string

const (
	SampleTypePretrain    SampleType = "pretrain"
	SampleTypeInstruct    SampleType = "instruct"
	SampleTypeBatchReq    SampleType = "batch_request"
	SampleTypeBatchResult SampleType = "batch_result"
	SampleTypeBatchError  SampleType = "batch_error"
)

// Source indicates how a file was created.
type Source string

const (
	SourceUpload     Source = "upload"
	SourceRepository Source = "repository"
	SourceMistral    Source = "mistral"
)

// File represents an uploaded file.
type File struct {
	ID         string     `json:"id"`
	Object     string     `json:"object"`
	Bytes      int        `json:"bytes"`
	CreatedAt  int64      `json:"created_at"`
	Filename   string     `json:"filename"`
	Purpose    Purpose    `json:"purpose"`
	SampleType SampleType `json:"sample_type"`
	NumLines   *int       `json:"num_lines,omitempty"`
	MimeType   *string    `json:"mimetype,omitempty"`
	Source     Source     `json:"source"`
	Signature  *string    `json:"signature,omitempty"`
	Deleted    bool       `json:"deleted,omitempty"`
}

// ListResponse is the response from listing files.
type ListResponse struct {
	Data   []File `json:"data"`
	Object string `json:"object"`
	Total  *int   `json:"total,omitempty"`
}

// DeleteResponse is the response from deleting a file.
type DeleteResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Deleted bool   `json:"deleted"`
}

// SignedURL is the response from getting a signed file URL.
type SignedURL struct {
	URL string `json:"url"`
}

// Visibility controls who can see a file.
type Visibility string

const (
	VisibilitySharedGlobal    Visibility = "shared_global"
	VisibilitySharedOrg       Visibility = "shared_org"
	VisibilitySharedWorkspace Visibility = "shared_workspace"
	VisibilityPrivate         Visibility = "private"
)

// UploadParams holds parameters for uploading a file.
type UploadParams struct {
	Purpose    Purpose
	Expiry     *int
	Visibility *Visibility
}

// ListParams holds optional parameters for listing files.
type ListParams struct {
	Page     *int
	PageSize *int
	Purpose  *Purpose
	Search   *string
}
