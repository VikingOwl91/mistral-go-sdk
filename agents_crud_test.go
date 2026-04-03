package mistral

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/VikingOwl91/mistral-go-sdk/agents"
)

func TestCreateAgent_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("got method %s", r.Method)
		}
		if r.URL.Path != "/v1/agents" {
			t.Errorf("got path %s", r.URL.Path)
		}
		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)
		if body["name"] != "TestAgent" {
			t.Errorf("got name %v", body["name"])
		}
		if body["model"] != "mistral-small-latest" {
			t.Errorf("got model %v", body["model"])
		}
		json.NewEncoder(w).Encode(map[string]any{
			"object": "agent", "id": "ag-123", "name": "TestAgent",
			"model": "mistral-small-latest", "version": 1,
			"versions": []int{1}, "created_at": "2024-01-01T00:00:00Z",
			"updated_at": "2024-01-01T00:00:00Z", "deployment_chat": false, "source": "api",
		})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	agent, err := client.CreateAgent(context.Background(), &agents.CreateRequest{
		Name:  "TestAgent",
		Model: "mistral-small-latest",
	})
	if err != nil {
		t.Fatal(err)
	}
	if agent.ID != "ag-123" {
		t.Errorf("got id %q", agent.ID)
	}
	if agent.Name != "TestAgent" {
		t.Errorf("got name %q", agent.Name)
	}
}

func TestGetAgent_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("got method %s", r.Method)
		}
		if r.URL.Path != "/v1/agents/ag-123" {
			t.Errorf("got path %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"object": "agent", "id": "ag-123", "name": "TestAgent",
			"model": "mistral-small-latest", "version": 1,
			"versions": []int{1}, "created_at": "2024-01-01T00:00:00Z",
			"updated_at": "2024-01-01T00:00:00Z", "deployment_chat": false, "source": "api",
		})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	agent, err := client.GetAgent(context.Background(), "ag-123")
	if err != nil {
		t.Fatal(err)
	}
	if agent.ID != "ag-123" {
		t.Errorf("got id %q", agent.ID)
	}
}

func TestListAgents_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("got method %s", r.Method)
		}
		if r.URL.Query().Get("page") != "0" {
			t.Errorf("got page %q", r.URL.Query().Get("page"))
		}
		json.NewEncoder(w).Encode([]map[string]any{
			{
				"object": "agent", "id": "ag-1", "name": "A1",
				"model": "m", "version": 1, "versions": []int{1},
				"created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z",
				"deployment_chat": false, "source": "api",
			},
		})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	page := 0
	list, err := client.ListAgents(context.Background(), &agents.ListParams{Page: &page})
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 {
		t.Fatalf("got %d agents", len(list))
	}
	if list[0].ID != "ag-1" {
		t.Errorf("got id %q", list[0].ID)
	}
}

func TestUpdateAgent_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PATCH" {
			t.Errorf("got method %s", r.Method)
		}
		if r.URL.Path != "/v1/agents/ag-123" {
			t.Errorf("got path %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"object": "agent", "id": "ag-123", "name": "Updated",
			"model": "mistral-small-latest", "version": 2,
			"versions": []int{1, 2}, "created_at": "2024-01-01T00:00:00Z",
			"updated_at": "2024-01-01T00:01:00Z", "deployment_chat": false, "source": "api",
		})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	name := "Updated"
	agent, err := client.UpdateAgent(context.Background(), "ag-123", &agents.UpdateRequest{
		Name: &name,
	})
	if err != nil {
		t.Fatal(err)
	}
	if agent.Version != 2 {
		t.Errorf("got version %d", agent.Version)
	}
}

func TestDeleteAgent_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("got method %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	err := client.DeleteAgent(context.Background(), "ag-123")
	if err != nil {
		t.Fatal(err)
	}
}

func TestUpdateAgentVersion_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PATCH" {
			t.Errorf("got method %s", r.Method)
		}
		if r.URL.Query().Get("version") != "2" {
			t.Errorf("got version %q", r.URL.Query().Get("version"))
		}
		json.NewEncoder(w).Encode(map[string]any{
			"object": "agent", "id": "ag-123", "name": "A",
			"model": "m", "version": 2, "versions": []int{1, 2},
			"created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z",
			"deployment_chat": false, "source": "api",
		})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	agent, err := client.UpdateAgentVersion(context.Background(), "ag-123", 2)
	if err != nil {
		t.Fatal(err)
	}
	if agent.Version != 2 {
		t.Errorf("got version %d", agent.Version)
	}
}

func TestListAgentVersions_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/agents/ag-123/versions" {
			t.Errorf("got path %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode([]map[string]any{
			{
				"object": "agent", "id": "ag-123", "name": "A",
				"model": "m", "version": 1, "versions": []int{1, 2},
				"created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z",
				"deployment_chat": false, "source": "api",
			},
			{
				"object": "agent", "id": "ag-123", "name": "A",
				"model": "m", "version": 2, "versions": []int{1, 2},
				"created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:01:00Z",
				"deployment_chat": false, "source": "api",
			},
		})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	versions, err := client.ListAgentVersions(context.Background(), "ag-123", nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(versions) != 2 {
		t.Fatalf("got %d versions", len(versions))
	}
}

func TestSetAgentAlias_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("got method %s", r.Method)
		}
		if r.URL.Query().Get("alias") != "production" {
			t.Errorf("got alias %q", r.URL.Query().Get("alias"))
		}
		if r.URL.Query().Get("version") != "2" {
			t.Errorf("got version %q", r.URL.Query().Get("version"))
		}
		json.NewEncoder(w).Encode(map[string]any{
			"alias": "production", "version": 2,
			"created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z",
		})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	alias, err := client.SetAgentAlias(context.Background(), "ag-123", "production", 2)
	if err != nil {
		t.Fatal(err)
	}
	if alias.Alias != "production" {
		t.Errorf("got alias %q", alias.Alias)
	}
	if alias.Version != 2 {
		t.Errorf("got version %d", alias.Version)
	}
}

func TestListAgentAliases_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/agents/ag-123/aliases" {
			t.Errorf("got path %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode([]map[string]any{
			{"alias": "production", "version": 2, "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"},
		})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	aliases, err := client.ListAgentAliases(context.Background(), "ag-123")
	if err != nil {
		t.Fatal(err)
	}
	if len(aliases) != 1 {
		t.Fatalf("got %d aliases", len(aliases))
	}
}

func TestDeleteAgentAlias_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("got method %s", r.Method)
		}
		if r.URL.Query().Get("alias") != "staging" {
			t.Errorf("got alias %q", r.URL.Query().Get("alias"))
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	err := client.DeleteAgentAlias(context.Background(), "ag-123", "staging")
	if err != nil {
		t.Fatal(err)
	}
}
