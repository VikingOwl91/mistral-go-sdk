package mistral

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/VikingOwl91/mistral-go-sdk/connector"
)

// CreateConnector registers a new MCP connector.
func (c *Client) CreateConnector(ctx context.Context, req *connector.CreateRequest) (*connector.Connector, error) {
	var resp connector.Connector
	if err := c.doJSON(ctx, "POST", "/v1/connectors", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ListConnectors returns all connectors.
func (c *Client) ListConnectors(ctx context.Context, params *connector.ListParams) ([]connector.Connector, error) {
	path := "/v1/connectors"
	if params != nil {
		q := url.Values{}
		if params.Page != nil {
			q.Set("page", strconv.Itoa(*params.Page))
		}
		if params.PageSize != nil {
			q.Set("page_size", strconv.Itoa(*params.PageSize))
		}
		if encoded := q.Encode(); encoded != "" {
			path += "?" + encoded
		}
	}
	var resp []connector.Connector
	if err := c.doJSON(ctx, "GET", path, nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// GetConnector retrieves a connector by ID or name.
func (c *Client) GetConnector(ctx context.Context, idOrName string) (*connector.Connector, error) {
	var resp connector.Connector
	if err := c.doJSON(ctx, "GET", fmt.Sprintf("/v1/connectors/%s", idOrName), nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// UpdateConnector updates an existing connector.
func (c *Client) UpdateConnector(ctx context.Context, idOrName string, req *connector.UpdateRequest) (*connector.Connector, error) {
	var resp connector.Connector
	if err := c.doJSON(ctx, "PATCH", fmt.Sprintf("/v1/connectors/%s", idOrName), req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// DeleteConnector deletes a connector.
func (c *Client) DeleteConnector(ctx context.Context, idOrName string) error {
	resp, err := c.do(ctx, "DELETE", fmt.Sprintf("/v1/connectors/%s", idOrName), nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return parseAPIError(resp)
	}
	return nil
}

// GetConnectorAuthURL returns the OAuth2 authorization URL for a connector.
func (c *Client) GetConnectorAuthURL(ctx context.Context, idOrName string, appReturnURL *string) (*connector.AuthURLResponse, error) {
	path := fmt.Sprintf("/v1/connectors/%s/auth_url", idOrName)
	if appReturnURL != nil {
		path += "?app_return_url=" + url.QueryEscape(*appReturnURL)
	}
	var resp connector.AuthURLResponse
	if err := c.doJSON(ctx, "GET", path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ListConnectorTools lists tools exposed by a connector.
func (c *Client) ListConnectorTools(ctx context.Context, idOrName string, params *connector.ListToolsParams) ([]connector.Tool, error) {
	path := fmt.Sprintf("/v1/connectors/%s/tools", idOrName)
	if params != nil {
		q := url.Values{}
		if params.Page != nil {
			q.Set("page", strconv.Itoa(*params.Page))
		}
		if params.PageSize != nil {
			q.Set("page_size", strconv.Itoa(*params.PageSize))
		}
		if params.Refresh != nil {
			q.Set("refresh", strconv.FormatBool(*params.Refresh))
		}
		if encoded := q.Encode(); encoded != "" {
			path += "?" + encoded
		}
	}
	var resp []connector.Tool
	if err := c.doJSON(ctx, "GET", path, nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// CallConnectorTool invokes a tool on a connector.
func (c *Client) CallConnectorTool(ctx context.Context, idOrName, toolName string, req *connector.CallToolRequest) (*connector.CallToolResponse, error) {
	var resp connector.CallToolResponse
	if err := c.doJSON(ctx, "POST", fmt.Sprintf("/v1/connectors/%s/tools/%s/call", idOrName, toolName), req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
