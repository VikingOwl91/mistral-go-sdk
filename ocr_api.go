package mistral

import (
	"context"

	"somegit.dev/vikingowl/mistral-go-sdk/ocr"
)

// OCR performs optical character recognition on a document.
func (c *Client) OCR(ctx context.Context, req *ocr.Request) (*ocr.Response, error) {
	var resp ocr.Response
	if err := c.doJSON(ctx, "POST", "/v1/ocr", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
