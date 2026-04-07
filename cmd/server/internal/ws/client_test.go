package ws

import (
	"context"
	"encoding/json"
	"testing"
)

// ---------------------------------------------------------------------------
// dispatchRoomMessage validation — no DB needed (returns before CreateMessage)
// ---------------------------------------------------------------------------

// Test 13 – room_id zero: no broadcast
func TestDispatchRoomMessage_InvalidRoomID(t *testing.T) {
	h := startHub(t)
	c := newTestClient(h, 1, map[int64]bool{10: true})
	sync := newTestClient(h, 99, map[int64]bool{})
	registerAll(t, h, c, sync)

	payload, _ := json.Marshal(RoomMessagePayload{RoomID: 0, Content: "hi"})
	c.dispatchRoomMessage(Message{Type: TypeRoomMessage, Payload: payload}, context.Background())
	syncHub(t, h, sync)

	expectNoMessage(t, c.send)
}

// Test 14 – empty content: no broadcast
func TestDispatchRoomMessage_EmptyContent(t *testing.T) {
	h := startHub(t)
	c := newTestClient(h, 1, map[int64]bool{10: true})
	sync := newTestClient(h, 99, map[int64]bool{})
	registerAll(t, h, c, sync)

	payload, _ := json.Marshal(RoomMessagePayload{RoomID: 10, Content: ""})
	c.dispatchRoomMessage(Message{Type: TypeRoomMessage, Payload: payload}, context.Background())
	syncHub(t, h, sync)

	expectNoMessage(t, c.send)
}

// Test 15 – client not in room: no broadcast
func TestDispatchRoomMessage_ClientNotInRoom(t *testing.T) {
	h := startHub(t)
	c := newTestClient(h, 1, map[int64]bool{10: true})
	sync := newTestClient(h, 99, map[int64]bool{})
	registerAll(t, h, c, sync)

	payload, _ := json.Marshal(RoomMessagePayload{RoomID: 99, Content: "hi"})
	c.dispatchRoomMessage(Message{Type: TypeRoomMessage, Payload: payload}, context.Background())
	syncHub(t, h, sync)

	expectNoMessage(t, c.send)
}

// Test 16 – malformed payload: no panic, no broadcast
func TestDispatchRoomMessage_MalformedPayload(t *testing.T) {
	h := startHub(t)
	c := newTestClient(h, 1, map[int64]bool{10: true})
	sync := newTestClient(h, 99, map[int64]bool{})
	registerAll(t, h, c, sync)

	c.dispatchRoomMessage(Message{Type: TypeRoomMessage, Payload: json.RawMessage(`not json`)}, context.Background())
	syncHub(t, h, sync)

	expectNoMessage(t, c.send)
}
