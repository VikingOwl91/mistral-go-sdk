package mistral

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"somegit.dev/vikingowl/mistral-go-sdk/connector"
)

func TestCreateConnector_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/v1/connectors" {
			t.Errorf("got path %s", r.URL.Path)
		}
		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)
		if body["name"] != "my_connector" {
			t.Errorf("got name %v", body["name"])
		}
		if body["server"] != "https://mcp.example.com" {
			t.Errorf("got server %v", body["server"])
		}
		json.NewEncoder(w).Encode(map[string]any{
			"id": "conn-1", "name": "my_connector",
			"description": "test", "created_at": "2025-01-01",
			"modified_at": "2025-01-01", "server": "https://mcp.example.com",
		})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	resp, err := client.CreateConnector(context.Background(), &connector.CreateRequest{
		Name:        "my_connector",
		Description: "test",
		Server:      "https://mcp.example.com",
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.ID != "conn-1" {
		t.Errorf("got id %q", resp.ID)
	}
}

func TestListConnectors_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("expected GET, got %s", r.Method)
		}
		json.NewEncoder(w).Encode([]map[string]any{
			{"id": "c1", "name": "conn1", "description": "d1", "created_at": "t", "modified_at": "t"},
		})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	list, err := client.ListConnectors(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 {
		t.Fatalf("got %d connectors", len(list))
	}
	if list[0].ID != "c1" {
		t.Errorf("got id %q", list[0].ID)
	}
}

func TestGetConnector_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/connectors/my_conn" {
			t.Errorf("got path %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"id": "c1", "name": "my_conn", "description": "d",
			"created_at": "t", "modified_at": "t",
		})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	c, err := client.GetConnector(context.Background(), "my_conn")
	if err != nil {
		t.Fatal(err)
	}
	if c.Name != "my_conn" {
		t.Errorf("got name %q", c.Name)
	}
}

func TestUpdateConnector_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PATCH" {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"id": "c1", "name": "updated", "description": "new desc",
			"created_at": "t", "modified_at": "t",
		})
	}))
	defer server.Close()

	name := "updated"
	client := NewClient("key", WithBaseURL(server.URL))
	resp, err := client.UpdateConnector(context.Background(), "c1", &connector.UpdateRequest{
		Name: &name,
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Name != "updated" {
		t.Errorf("got name %q", resp.Name)
	}
}

func TestDeleteConnector_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	err := client.DeleteConnector(context.Background(), "c1")
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetConnectorAuthURL_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/connectors/c1/auth_url" {
			t.Errorf("got path %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"auth_url": "https://oauth.example.com/authorize",
			"ttl":      3600,
		})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	resp, err := client.GetConnectorAuthURL(context.Background(), "c1", nil)
	if err != nil {
		t.Fatal(err)
	}
	if resp.AuthURL != "https://oauth.example.com/authorize" {
		t.Errorf("got auth_url %q", resp.AuthURL)
	}
	if resp.TTL != 3600 {
		t.Errorf("got ttl %d", resp.TTL)
	}
}

func TestListConnectorTools_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/connectors/c1/tools" {
			t.Errorf("got path %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode([]map[string]any{
			{"id": "t1", "name": "search", "description": "search the web"},
		})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	tools, err := client.ListConnectorTools(context.Background(), "c1", nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(tools) != 1 {
		t.Fatalf("got %d tools", len(tools))
	}
	if tools[0].Name != "search" {
		t.Errorf("got name %q", tools[0].Name)
	}
}

func TestCallConnectorTool_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/v1/connectors/c1/tools/search/call" {
			t.Errorf("got path %s", r.URL.Path)
		}
		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)
		args := body["arguments"].(map[string]any)
		if args["query"] != "hello" {
			t.Errorf("got query %v", args["query"])
		}
		json.NewEncoder(w).Encode(map[string]any{
			"content": []map[string]any{{"type": "text", "text": "result"}},
		})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	resp, err := client.CallConnectorTool(context.Background(), "c1", "search", &connector.CallToolRequest{
		Arguments: map[string]any{"query": "hello"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Content == nil {
		t.Error("expected non-nil content")
	}
}
