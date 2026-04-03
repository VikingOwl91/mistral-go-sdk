package mistral

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/VikingOwl91/mistral-go-sdk/agents"
)

// CreateAgent creates a new agent.
func (c *Client) CreateAgent(ctx context.Context, req *agents.CreateRequest) (*agents.Agent, error) {
	var resp agents.Agent
	if err := c.doJSON(ctx, "POST", "/v1/agents", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ListAgents lists agents with optional filters and pagination.
func (c *Client) ListAgents(ctx context.Context, params *agents.ListParams) ([]agents.Agent, error) {
	path := "/v1/agents"
	if params != nil {
		q := url.Values{}
		if params.Page != nil {
			q.Set("page", strconv.Itoa(*params.Page))
		}
		if params.PageSize != nil {
			q.Set("page_size", strconv.Itoa(*params.PageSize))
		}
		if params.DeploymentChat != nil {
			q.Set("deployment_chat", strconv.FormatBool(*params.DeploymentChat))
		}
		if params.Name != nil {
			q.Set("name", *params.Name)
		}
		if params.Search != nil {
			q.Set("search", *params.Search)
		}
		if params.ID != nil {
			q.Set("id", *params.ID)
		}
		if encoded := q.Encode(); encoded != "" {
			path += "?" + encoded
		}
	}
	var resp []agents.Agent
	if err := c.doJSON(ctx, "GET", path, nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// GetAgent retrieves an agent by ID.
func (c *Client) GetAgent(ctx context.Context, agentID string) (*agents.Agent, error) {
	var resp agents.Agent
	if err := c.doJSON(ctx, "GET", fmt.Sprintf("/v1/agents/%s", agentID), nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// UpdateAgent updates an agent, creating a new version.
func (c *Client) UpdateAgent(ctx context.Context, agentID string, req *agents.UpdateRequest) (*agents.Agent, error) {
	var resp agents.Agent
	if err := c.doJSON(ctx, "PATCH", fmt.Sprintf("/v1/agents/%s", agentID), req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// DeleteAgent deletes an agent.
func (c *Client) DeleteAgent(ctx context.Context, agentID string) error {
	resp, err := c.do(ctx, "DELETE", fmt.Sprintf("/v1/agents/%s", agentID), nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return parseAPIError(resp)
	}
	return nil
}

// UpdateAgentVersion switches an agent to a specific version.
func (c *Client) UpdateAgentVersion(ctx context.Context, agentID string, version int) (*agents.Agent, error) {
	path := fmt.Sprintf("/v1/agents/%s/version?version=%d", agentID, version)
	var resp agents.Agent
	if err := c.doJSON(ctx, "PATCH", path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ListAgentVersions lists all versions of an agent.
func (c *Client) ListAgentVersions(ctx context.Context, agentID string, params *agents.ListVersionsParams) ([]agents.Agent, error) {
	path := fmt.Sprintf("/v1/agents/%s/versions", agentID)
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
	var resp []agents.Agent
	if err := c.doJSON(ctx, "GET", path, nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// GetAgentVersion retrieves a specific version of an agent.
func (c *Client) GetAgentVersion(ctx context.Context, agentID string, version string) (*agents.Agent, error) {
	var resp agents.Agent
	if err := c.doJSON(ctx, "GET", fmt.Sprintf("/v1/agents/%s/versions/%s", agentID, version), nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// SetAgentAlias creates or updates an agent version alias.
func (c *Client) SetAgentAlias(ctx context.Context, agentID string, alias string, version int) (*agents.AliasResponse, error) {
	path := fmt.Sprintf("/v1/agents/%s/aliases?alias=%s&version=%d", agentID, url.QueryEscape(alias), version)
	var resp agents.AliasResponse
	if err := c.doJSON(ctx, "PUT", path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ListAgentAliases lists all version aliases for an agent.
func (c *Client) ListAgentAliases(ctx context.Context, agentID string) ([]agents.AliasResponse, error) {
	var resp []agents.AliasResponse
	if err := c.doJSON(ctx, "GET", fmt.Sprintf("/v1/agents/%s/aliases", agentID), nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// DeleteAgentAlias deletes an agent version alias.
func (c *Client) DeleteAgentAlias(ctx context.Context, agentID string, alias string) error {
	path := fmt.Sprintf("/v1/agents/%s/aliases?alias=%s", agentID, url.QueryEscape(alias))
	resp, err := c.do(ctx, "DELETE", path, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return parseAPIError(resp)
	}
	return nil
}
