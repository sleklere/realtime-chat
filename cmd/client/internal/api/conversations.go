package api

import "fmt"

// ListConversations returns all DM conversations for the authenticated user.
func (c *Client) ListConversations() ([]ConversationResponse, error) {
	var convs []ConversationResponse
	if err := c.do("GET", "/api/v1/conversations", nil, &convs); err != nil {
		return nil, err
	}
	return convs, nil
}

// GetConversationMessages returns the message history for a conversation.
func (c *Client) GetConversationMessages(conversationID int64, limit int) ([]MessageResponse, error) {
	var msgs []MessageResponse
	path := fmt.Sprintf("/api/v1/conversations/%d/messages?limit=%d", conversationID, limit)
	if err := c.do("GET", path, nil, &msgs); err != nil {
		return nil, err
	}
	return msgs, nil
}
