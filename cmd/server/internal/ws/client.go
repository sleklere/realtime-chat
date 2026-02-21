package ws

import (
	"context"
	"encoding/json"
)

func (c *Client) ReadPump(ctx context.Context) {
	defer func() {
		c.hub.unregister <- c
	}()

	for {
		// read from websocket
		_, data, err := c.conn.Read(ctx)
		if err != nil {
			c.logger.Debug("ws read error", "user_id", c.userID, "err", err)
			return
		}
		// parse envelope
		var msg Message
		err = json.Unmarshal(data, &msg)
		if err != nil {
			c.logger.Warn("ws invalid message json", "user_id", c.userID, "err", err)
			continue
		}
		// route by message type
		switch msg.Type {
		case TypeRoomMessage:
			// parse room message payload
			var roomMsgPayload RoomMessagePayload
			err = json.Unmarshal(msg.Payload, &roomMsgPayload)
			if err != nil {
				c.logger.Warn("ws invalid payload json", "user_id", c.userID, "type", msg.Type, "err", err)
				continue
			}
			// validate: roomID > 0, content not empty, client is a room member
			_, clientInRoom := c.roomIDs[roomMsgPayload.RoomID]
			if roomMsgPayload.RoomID == 0 || roomMsgPayload.Content == "" ||
				!clientInRoom {
				c.logger.Warn("ws room message validation failed", "user_id", c.userID, "room_id", roomMsgPayload.RoomID)
				continue
			}
			// populate server fields and broadcast
			roomMsgPayload.SenderID = c.userID
			roomMsgPayload.SenderUsername = c.username
			completePayload, err := json.Marshal(roomMsgPayload)
			if err != nil {
				c.logger.Warn("ws payload marshal error", "user_id", c.userID, "err", err)
				continue
			}
			msgWithCompletePayload := Message{Type: msg.Type, Payload: completePayload, Timestamp: msg.Timestamp}

			broadcastMsg := BroadcastMsg{msg: msgWithCompletePayload, targetRoomID: roomMsgPayload.RoomID}
			c.hub.broadcast <- broadcastMsg
		}
	}
}
