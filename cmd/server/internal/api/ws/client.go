package ws

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// time to write to the socket
	writeWait = 10 * time.Second

	// time to read pong response
	pongWait = 60 * time.Second

	// send ping every now and then (must be less than pongWait)
	pingPeriod = (pongWait * 9) / 10

	maxMessageSize = 512
)

type Client struct {
	conn   *websocket.Conn
	send   chan []byte
	userID int32
}

func (c *Client) readPump() {
	defer func() {
		c.conn.Close()
	}()

	// limits and timeouts
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		// read msg from client
		_, data, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err,
				websocket.CloseGoingAway,
				websocket.CloseAbnormalClosure) {
				log.Printf("error websocket: %v", err)
			}
			break
		}

		// parse msg
		var msg Message
		if err := json.Unmarshal(data, &msg); err != nil {
			log.Printf("error parsing message: %v", err)
			continue
		}

		// just log for now
		log.Printf("message recieved from user %d: tipo=%s", c.userID, msg.Type)

		// echo: re-send the message (temporary, testing)
		c.send <- data
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// chann closed
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			// write msg as text
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Printf("error creating writer: %v", err)
				return
			}
			w.Write(message)

			// send all pending messages in the buffer
			n := len(c.send)
			for range n {
				w.Write([]byte{'\n'})
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				log.Printf("error closing writer: %v", err)
				return
			}
		case <-ticker.C:
			// send ping periodically
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("error writing ping message: %v", err)
				return
			}

		}
	}
}
