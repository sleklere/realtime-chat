package ws

import (
	"context"
	"encoding/json"
	"testing"
)

// ---------------------------------------------------------------------------
// Test 13 – handleMessage TypeJoinRoom: Hub updates state, client receives
//           broadcasts in the new room
// ---------------------------------------------------------------------------

func TestHandleMessage_JoinRoom(t *testing.T) {
	h := startHub(t)

	c := newTestClient(h, 1, map[int64]bool{})
	sync := newTestClient(h, 99, map[int64]bool{})
	registerAll(t, h, c, sync)

	payload, _ := json.Marshal(JoinRoomPayload{RoomID: 5})
	c.dispatchWSMessage(Message{Type: TypeJoinRoom, Payload: payload}, context.Background())
	syncHub(t, h, sync)

	// broadcast to room 5 — c should now receive
	h.broadcast <- BroadcastMsg{msg: Message{Type: TypeRoomMessage}, targetRoomID: 5}
	syncHub(t, h, sync)

	got := expectMessage(t, c.send)
	if got.Type != TypeRoomMessage {
		t.Fatalf("expected %s, got %s", TypeRoomMessage, got.Type)
	}
}

// ---------------------------------------------------------------------------
// Test 14 – handleMessage TypeLeaveRoom: client stops receiving
// ---------------------------------------------------------------------------

func TestHandleMessage_LeaveRoom(t *testing.T) {
	h := startHub(t)

	c := newTestClient(h, 1, map[int64]bool{5: true})
	sync := newTestClient(h, 99, map[int64]bool{})
	registerAll(t, h, c, sync)

	// precondition: c receives in room 5
	h.broadcast <- BroadcastMsg{msg: Message{Type: TypeRoomMessage}, targetRoomID: 5}
	syncHub(t, h, sync)
	pre := expectMessage(t, c.send)
	if pre.Type != TypeRoomMessage {
		t.Fatalf("precondition failed: expected %s, got %s", TypeRoomMessage, pre.Type)
	}

	// leave room 5
	payload, _ := json.Marshal(JoinRoomPayload{RoomID: 5})
	c.dispatchWSMessage(Message{Type: TypeLeaveRoom, Payload: payload}, context.Background())
	syncHub(t, h, sync)

	// broadcast again — c should NOT receive
	h.broadcast <- BroadcastMsg{msg: Message{Type: TypeRoomMessage}, targetRoomID: 5}
	syncHub(t, h, sync)

	expectNoMessage(t, c.send)
}

// ---------------------------------------------------------------------------
// Test 15 – handleMessage TypeDirectMessage: sender and recipient both receive
// ---------------------------------------------------------------------------

func TestHandleMessage_DirectMessage(t *testing.T) {
	h := startHub(t)

	a := newTestClient(h, 1, map[int64]bool{})
	b := newTestClient(h, 2, map[int64]bool{})
	sync := newTestClient(h, 99, map[int64]bool{})
	registerAll(t, h, a, b, sync)

	payload, _ := json.Marshal(DirectMessagePayload{ToUserID: 2, Content: "hello"})
	a.dispatchWSMessage(Message{Type: TypeDirectMessage, Payload: payload}, context.Background())
	syncHub(t, h, sync)

	// sender (a) gets a copy
	gotA := expectMessage(t, a.send)
	if gotA.Type != TypeDirectMessage {
		t.Fatalf("sender: expected %s, got %s", TypeDirectMessage, gotA.Type)
	}

	// recipient (b) gets the message
	gotB := expectMessage(t, b.send)
	if gotB.Type != TypeDirectMessage {
		t.Fatalf("recipient: expected %s, got %s", TypeDirectMessage, gotB.Type)
	}

	// verify server-populated fields
	var dm DirectMessagePayload
	if err := json.Unmarshal(gotB.Payload, &dm); err != nil {
		t.Fatal(err)
	}
	if dm.FromUserID != 1 {
		t.Fatalf("expected from_user_id=1, got %d", dm.FromUserID)
	}
	if dm.FromUsername != "user_1" {
		t.Fatalf("expected from_username=user_1, got %s", dm.FromUsername)
	}
	if dm.Content != "hello" {
		t.Fatalf("expected content='hello', got '%s'", dm.Content)
	}
}

// ---------------------------------------------------------------------------
// Test 16 – handleMessage unknown type: no panic, no message sent
// ---------------------------------------------------------------------------

func TestHandleMessage_UnknownType(t *testing.T) {
	h := startHub(t)

	c := newTestClient(h, 1, map[int64]bool{})
	sync := newTestClient(h, 99, map[int64]bool{})
	registerAll(t, h, c, sync)

	c.dispatchWSMessage(Message{
		Type:    "xxx",
		Payload: json.RawMessage(`{}`),
	}, context.Background())

	expectNoMessage(t, c.send)
}
