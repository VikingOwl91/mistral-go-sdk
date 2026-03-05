package mistral

import (
	"context"
	"encoding/json"
	"fmt"

	"somegit.dev/vikingowl/mistral-go-sdk/conversation"
)

// StartConversation creates and starts a new conversation.
func (c *Client) StartConversation(ctx context.Context, req *conversation.StartRequest) (*conversation.Response, error) {
	var resp conversation.Response
	if err := c.doJSON(ctx, "POST", "/v1/conversations", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// StartConversationStream creates a conversation and returns a stream of events.
func (c *Client) StartConversationStream(ctx context.Context, req *conversation.StartRequest) (*EventStream, error) {
	req.SetStream(true)
	resp, err := c.doStream(ctx, "POST", "/v1/conversations", req)
	if err != nil {
		return nil, err
	}
	return newEventStream(resp.Body), nil
}

// AppendConversation appends inputs to an existing conversation.
func (c *Client) AppendConversation(ctx context.Context, conversationID string, req *conversation.AppendRequest) (*conversation.Response, error) {
	var resp conversation.Response
	if err := c.doJSON(ctx, "POST", fmt.Sprintf("/v1/conversations/%s", conversationID), req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// AppendConversationStream appends to a conversation and returns a stream of events.
func (c *Client) AppendConversationStream(ctx context.Context, conversationID string, req *conversation.AppendRequest) (*EventStream, error) {
	req.SetStream(true)
	resp, err := c.doStream(ctx, "POST", fmt.Sprintf("/v1/conversations/%s", conversationID), req)
	if err != nil {
		return nil, err
	}
	return newEventStream(resp.Body), nil
}

// RestartConversation restarts a conversation from a specific entry.
func (c *Client) RestartConversation(ctx context.Context, conversationID string, req *conversation.RestartRequest) (*conversation.Response, error) {
	var resp conversation.Response
	if err := c.doJSON(ctx, "POST", fmt.Sprintf("/v1/conversations/%s/restart", conversationID), req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// RestartConversationStream restarts a conversation and returns a stream of events.
func (c *Client) RestartConversationStream(ctx context.Context, conversationID string, req *conversation.RestartRequest) (*EventStream, error) {
	req.SetStream(true)
	resp, err := c.doStream(ctx, "POST", fmt.Sprintf("/v1/conversations/%s/restart", conversationID), req)
	if err != nil {
		return nil, err
	}
	return newEventStream(resp.Body), nil
}

// GetConversation retrieves conversation metadata.
func (c *Client) GetConversation(ctx context.Context, conversationID string) (*conversation.Conversation, error) {
	var resp conversation.Conversation
	if err := c.doJSON(ctx, "GET", fmt.Sprintf("/v1/conversations/%s", conversationID), nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ListConversations lists conversations with optional pagination.
func (c *Client) ListConversations(ctx context.Context) ([]conversation.Conversation, error) {
	var resp []conversation.Conversation
	if err := c.doJSON(ctx, "GET", "/v1/conversations", nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// DeleteConversation deletes a conversation.
func (c *Client) DeleteConversation(ctx context.Context, conversationID string) error {
	resp, err := c.do(ctx, "DELETE", fmt.Sprintf("/v1/conversations/%s", conversationID), nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return parseAPIError(resp)
	}
	return nil
}

// GetConversationHistory returns the full history of a conversation.
func (c *Client) GetConversationHistory(ctx context.Context, conversationID string) (*conversation.History, error) {
	var resp conversation.History
	if err := c.doJSON(ctx, "GET", fmt.Sprintf("/v1/conversations/%s/history", conversationID), nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetConversationMessages returns the messages of a conversation.
func (c *Client) GetConversationMessages(ctx context.Context, conversationID string) (*conversation.Messages, error) {
	var resp conversation.Messages
	if err := c.doJSON(ctx, "GET", fmt.Sprintf("/v1/conversations/%s/messages", conversationID), nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// EventStream wraps the generic Stream to provide typed conversation events.
type EventStream struct {
	stream *Stream[json.RawMessage]
	event  conversation.Event
	err    error
}

func newEventStream(body readCloser) *EventStream {
	return &EventStream{
		stream: newStream[json.RawMessage](body),
	}
}

// Next advances to the next event. Returns false when done or on error.
func (s *EventStream) Next() bool {
	if s.err != nil {
		return false
	}
	if !s.stream.Next() {
		s.err = s.stream.Err()
		return false
	}
	event, err := conversation.UnmarshalEvent(s.stream.Current())
	if err != nil {
		s.err = err
		return false
	}
	s.event = event
	return true
}

// Current returns the most recently read event.
func (s *EventStream) Current() conversation.Event { return s.event }

// Err returns any error encountered during streaming.
func (s *EventStream) Err() error { return s.err }

// Close releases the underlying connection.
func (s *EventStream) Close() error { return s.stream.Close() }

type readCloser = interface {
	Read(p []byte) (n int, err error)
	Close() error
}
