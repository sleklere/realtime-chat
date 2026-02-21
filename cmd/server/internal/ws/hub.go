package ws

import (
	"github.com/coder/websocket"
)

// Client represents a single WebSocket connection.
type Client struct {
	hub      *Hub
	conn     *websocket.Conn
	userID   int64
	username string
	roomIDs  map[int64]bool // rooms this client is a member of
	send     chan Message   // Hub writes here, WritePump drains
}

// Hub is the central message broker.
// A single goroutine runs Hub.Run() and is the sole owner of the maps.
type Hub struct {
	clients map[int64]*Client        // userID → client
	rooms   map[int64]map[int64]bool // roomID → set of userIDs online

	register   chan *Client
	unregister chan *Client
	broadcast  chan BroadcastMsg
}

// BroadcastMsg wraps a message with routing info.
type BroadcastMsg struct {
	msg           Message
	targetRoomID  int64   // if > 0, route to room members
	targetUserIDs []int64 // if targetRoomID == 0, route to these users (DM)
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[int64]*Client),
		rooms:      make(map[int64]map[int64]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan BroadcastMsg, 256),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			for roomId, isInRoom := range client.roomIDs {
				_, ok := h.rooms[roomId]
				if !ok {
					h.rooms[roomId] = make(map[int64]bool)
				}
				h.rooms[roomId][client.userID] = isInRoom
			}
			h.clients[client.userID] = client
		case client := <-h.unregister:
			for roomId := range client.roomIDs {
				_, ok := h.rooms[roomId]
				if !ok {
					continue
				}
				delete(h.rooms[roomId], client.userID)
			}
			delete(h.clients, client.userID)
			close(client.send)
		case broadcastMsg := <-h.broadcast:
			// if it's a room msg
			if broadcastMsg.targetRoomID > 0 {
				for userId, isConnected := range h.rooms[broadcastMsg.targetRoomID] {
					if !isConnected {
						continue
					}
					c := h.clients[userId]
					select {
					case c.send <- broadcastMsg.msg:
					default:
						// client laggeado, lo sacamos
						h.kickClient(c)
					}
				}
			} else {
				for _, userId := range broadcastMsg.targetUserIDs {
					// if the client exists
					if c, ok := h.clients[userId]; ok {
						select {
						case c.send <- broadcastMsg.msg:
						default:
							// client laggeado, lo sacamos
							h.kickClient(c)
						}

					}
				}
			}
		}
	}
}

func (h *Hub) kickClient(c *Client) {
	close(c.send)
	delete(h.clients, c.userID)
	for roomId := range c.roomIDs {
		_, ok := h.rooms[roomId]
		if !ok {
			continue
		}
		delete(h.rooms[roomId], c.userID)
	}
}
