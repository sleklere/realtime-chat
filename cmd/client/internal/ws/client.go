// Package ws provides a WebSocket client for real-time chat messaging.
package ws

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/coder/websocket"
)

// IncomingMsg is a Bubble Tea message carrying a received WebSocket message.
type IncomingMsg struct {
	Message Message
}

// ErrorMsg is a Bubble Tea message carrying a WebSocket error.
type ErrorMsg struct {
	Err error
}

// ConnectedMsg signals that the WebSocket connection is established.
type ConnectedMsg struct{}

// Client manages a WebSocket connection to the chat server.
type Client struct {
	conn    *websocket.Conn
	sendCh  chan Message
	program *tea.Program
	logger  *slog.Logger
	cancel  context.CancelFunc
}

// Connect establishes a WebSocket connection and starts read/write loops.
func Connect(ctx context.Context, wsURL, token string, program *tea.Program, logger *slog.Logger) (*Client, error) {
	url := wsURL + "/api/v1/ws"

	headers := http.Header{}
	headers.Set("Authorization", "Bearer "+token)

	conn, _, err := websocket.Dial(ctx, url, &websocket.DialOptions{
		HTTPHeader: headers,
	})
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(ctx)

	c := &Client{
		conn:    conn,
		sendCh:  make(chan Message, 64),
		program: program,
		logger:  logger,
		cancel:  cancel,
	}

	go c.readLoop(ctx)
	go c.writeLoop(ctx)

	logger.Info("websocket connected", "url", url)
	return c, nil
}

// Send enqueues a message for sending over the WebSocket connection.
func (c *Client) Send(msg Message) {
	msg.Timestamp = time.Now()
	c.sendCh <- msg
}

// SendDirectMessage sends a direct message to the specified user.
func (c *Client) SendDirectMessage(toUserID int64, content string) error {
	payload, err := json.Marshal(DirectMessagePayload{
		ToUserID: toUserID,
		Content:  content,
	})
	if err != nil {
		return err
	}

	c.Send(Message{
		Type:    TypeDirectMessage,
		Payload: payload,
	})
	return nil
}

// SendRoomMessage sends a message to the specified room.
func (c *Client) SendRoomMessage(roomID int64, content string) error {
	payload, err := json.Marshal(RoomMessagePayload{
		RoomID:  roomID,
		Content: content,
	})
	if err != nil {
		return err
	}

	c.Send(Message{
		Type:    TypeRoomMessage,
		Payload: payload,
	})
	return nil
}

// Close cancels the connection context and closes the WebSocket.
func (c *Client) Close() {
	c.cancel()
	close(c.sendCh)
	_ = c.conn.Close(websocket.StatusNormalClosure, "bye")
	c.logger.Info("websocket closed")
}

func (c *Client) readLoop(ctx context.Context) {
	for {
		_, data, err := c.conn.Read(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			c.logger.Error("ws read error", "error", err)
			c.program.Send(ErrorMsg{Err: err})
			return
		}

		var msg Message
		if err := json.Unmarshal(data, &msg); err != nil {
			c.logger.Error("ws unmarshal error", "error", err)
			continue
		}

		c.logger.Debug("ws received", "type", msg.Type)
		c.program.Send(IncomingMsg{Message: msg})
	}
}

func (c *Client) writeLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-c.sendCh:
			if !ok {
				return
			}

			data, err := json.Marshal(msg)
			if err != nil {
				c.logger.Error("ws marshal error", "error", err)
				continue
			}

			if err := c.conn.Write(ctx, websocket.MessageText, data); err != nil {
				if ctx.Err() != nil {
					return
				}
				c.logger.Error("ws write error", "error", err)
				c.program.Send(ErrorMsg{Err: err})
				return
			}

			c.logger.Debug("ws sent", "type", msg.Type)
		}
	}
}
