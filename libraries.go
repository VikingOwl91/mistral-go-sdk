package mistral

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"strconv"

	"somegit.dev/vikingowl/mistral-go-sdk/library"
)

// CreateLibrary creates a new document library.
func (c *Client) CreateLibrary(ctx context.Context, req *library.CreateRequest) (*library.Library, error) {
	var resp library.Library
	if err := c.doJSON(ctx, "POST", "/v1/libraries", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ListLibraries lists all accessible libraries.
func (c *Client) ListLibraries(ctx context.Context) (*library.ListLibraryOut, error) {
	var resp library.ListLibraryOut
	if err := c.doJSON(ctx, "GET", "/v1/libraries", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetLibrary retrieves a library by ID.
func (c *Client) GetLibrary(ctx context.Context, libraryID string) (*library.Library, error) {
	var resp library.Library
	if err := c.doJSON(ctx, "GET", fmt.Sprintf("/v1/libraries/%s", libraryID), nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// UpdateLibrary updates a library's name and description.
func (c *Client) UpdateLibrary(ctx context.Context, libraryID string, req *library.UpdateRequest) (*library.Library, error) {
	var resp library.Library
	if err := c.doJSON(ctx, "PUT", fmt.Sprintf("/v1/libraries/%s", libraryID), req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// DeleteLibrary deletes a library and all its documents.
func (c *Client) DeleteLibrary(ctx context.Context, libraryID string) (*library.Library, error) {
	var resp library.Library
	if err := c.doJSON(ctx, "DELETE", fmt.Sprintf("/v1/libraries/%s", libraryID), nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// UploadDocument uploads a document to a library.
func (c *Client) UploadDocument(ctx context.Context, libraryID string, filename string, file io.Reader) (*library.Document, error) {
	var resp library.Document
	if err := c.doMultipart(ctx, fmt.Sprintf("/v1/libraries/%s/documents", libraryID), filename, file, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ListDocuments lists documents in a library.
func (c *Client) ListDocuments(ctx context.Context, libraryID string, params *library.ListDocumentParams) (*library.ListDocumentOut, error) {
	path := fmt.Sprintf("/v1/libraries/%s/documents", libraryID)
	if params != nil {
		q := url.Values{}
		if params.Search != nil {
			q.Set("search", *params.Search)
		}
		if params.PageSize != nil {
			q.Set("page_size", strconv.Itoa(*params.PageSize))
		}
		if params.Page != nil {
			q.Set("page", strconv.Itoa(*params.Page))
		}
		if params.SortBy != nil {
			q.Set("sort_by", *params.SortBy)
		}
		if params.SortOrder != nil {
			q.Set("sort_order", *params.SortOrder)
		}
		if encoded := q.Encode(); encoded != "" {
			path += "?" + encoded
		}
	}
	var resp library.ListDocumentOut
	if err := c.doJSON(ctx, "GET", path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetDocument retrieves a document's metadata.
func (c *Client) GetDocument(ctx context.Context, libraryID, documentID string) (*library.Document, error) {
	var resp library.Document
	if err := c.doJSON(ctx, "GET", fmt.Sprintf("/v1/libraries/%s/documents/%s", libraryID, documentID), nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// UpdateDocument updates a document's metadata.
func (c *Client) UpdateDocument(ctx context.Context, libraryID, documentID string, req *library.DocumentUpdateRequest) (*library.Document, error) {
	var resp library.Document
	if err := c.doJSON(ctx, "PUT", fmt.Sprintf("/v1/libraries/%s/documents/%s", libraryID, documentID), req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// DeleteDocument deletes a document from a library.
func (c *Client) DeleteDocument(ctx context.Context, libraryID, documentID string) error {
	resp, err := c.do(ctx, "DELETE", fmt.Sprintf("/v1/libraries/%s/documents/%s", libraryID, documentID), nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return parseAPIError(resp)
	}
	return nil
}

// GetDocumentTextContent retrieves the extracted text of a document.
func (c *Client) GetDocumentTextContent(ctx context.Context, libraryID, documentID string) (*library.DocumentTextContent, error) {
	var resp library.DocumentTextContent
	if err := c.doJSON(ctx, "GET", fmt.Sprintf("/v1/libraries/%s/documents/%s/text_content", libraryID, documentID), nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetDocumentStatus retrieves the processing status of a document.
func (c *Client) GetDocumentStatus(ctx context.Context, libraryID, documentID string) (*library.ProcessingStatusOut, error) {
	var resp library.ProcessingStatusOut
	if err := c.doJSON(ctx, "GET", fmt.Sprintf("/v1/libraries/%s/documents/%s/status", libraryID, documentID), nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetDocumentSignedURL retrieves a signed URL for downloading a document.
func (c *Client) GetDocumentSignedURL(ctx context.Context, libraryID, documentID string) (string, error) {
	var resp string
	if err := c.doJSON(ctx, "GET", fmt.Sprintf("/v1/libraries/%s/documents/%s/signed-url", libraryID, documentID), nil, &resp); err != nil {
		return "", err
	}
	return resp, nil
}

// GetDocumentExtractedTextSignedURL retrieves a signed URL for the extracted text.
func (c *Client) GetDocumentExtractedTextSignedURL(ctx context.Context, libraryID, documentID string) (string, error) {
	var resp string
	if err := c.doJSON(ctx, "GET", fmt.Sprintf("/v1/libraries/%s/documents/%s/extracted-text-signed-url", libraryID, documentID), nil, &resp); err != nil {
		return "", err
	}
	return resp, nil
}

// ReprocessDocument triggers reprocessing of a document.
func (c *Client) ReprocessDocument(ctx context.Context, libraryID, documentID string) error {
	resp, err := c.do(ctx, "POST", fmt.Sprintf("/v1/libraries/%s/documents/%s/reprocess", libraryID, documentID), nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return parseAPIError(resp)
	}
	return nil
}

// ListLibrarySharing lists all sharing entries for a library.
func (c *Client) ListLibrarySharing(ctx context.Context, libraryID string) (*library.ListSharingOut, error) {
	var resp library.ListSharingOut
	if err := c.doJSON(ctx, "GET", fmt.Sprintf("/v1/libraries/%s/share", libraryID), nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ShareLibrary creates or updates a sharing entry.
func (c *Client) ShareLibrary(ctx context.Context, libraryID string, req *library.SharingRequest) (*library.SharingOut, error) {
	var resp library.SharingOut
	if err := c.doJSON(ctx, "PUT", fmt.Sprintf("/v1/libraries/%s/share", libraryID), req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// UnshareLibrary deletes a sharing entry.
func (c *Client) UnshareLibrary(ctx context.Context, libraryID string, req *library.SharingDeleteRequest) (*library.SharingOut, error) {
	var resp library.SharingOut
	if err := c.doJSON(ctx, "DELETE", fmt.Sprintf("/v1/libraries/%s/share", libraryID), req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
