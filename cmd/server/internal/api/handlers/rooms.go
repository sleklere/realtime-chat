package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"
	reqdto "github.com/sleklere/realtime-chat/cmd/server/internal/api/dto/request"
	resdto "github.com/sleklere/realtime-chat/cmd/server/internal/api/dto/response"
	"github.com/sleklere/realtime-chat/cmd/server/internal/auth"
	"github.com/sleklere/realtime-chat/cmd/server/internal/httpx"
	dbstore "github.com/sleklere/realtime-chat/cmd/server/internal/store"
)

const (
	defaultMessageLimit = 50
	maxMessageLimit     = 100
)

// RoomHandler handles room-related HTTP requests.
type RoomHandler struct {
	queries *dbstore.Queries
	logger  *slog.Logger
}

// NewRoomHandler creates a new RoomHandler with the given queries and logger.
func NewRoomHandler(q *dbstore.Queries, l *slog.Logger) *RoomHandler {
	return &RoomHandler{queries: q, logger: l}
}

// Create handles room creation requests.
func (h *RoomHandler) Create(w http.ResponseWriter, r *http.Request) error {
	claims, ok := auth.ClaimsFromCtx(r.Context())
	if !ok {
		return httpx.New(http.StatusUnauthorized, "unauthorized", "missing claims", nil)
	}
	_ = claims

	var req reqdto.CreateRoomReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return httpx.BadRequest("invalid_json", "invalid json", err)
	}
	if req.Name == "" {
		return httpx.BadRequest("missing_name", "room name is required", nil)
	}

	slug := slugify(req.Name)
	room, err := h.queries.CreateRoom(r.Context(), dbstore.CreateRoomParams{
		Name: req.Name,
		Slug: slug,
	})
	if err != nil {
		return err
	}

	return httpx.JSON(w, http.StatusCreated, resdto.RoomRes{
		ID:        room.ID,
		Name:      room.Name,
		Slug:      room.Slug,
		CreatedAt: room.CreatedAt.Time,
	})
}

// List handles listing all rooms.
func (h *RoomHandler) List(w http.ResponseWriter, r *http.Request) error {
	rooms, err := h.queries.ListRooms(r.Context())
	if err != nil {
		return err
	}

	res := make([]resdto.RoomRes, len(rooms))
	for i, room := range rooms {
		res[i] = resdto.RoomRes{
			ID:        room.ID,
			Name:      room.Name,
			Slug:      room.Slug,
			CreatedAt: room.CreatedAt.Time,
		}
	}
	return httpx.JSON(w, http.StatusOK, res)
}

// GetBySlug handles fetching a room by its slug.
func (h *RoomHandler) GetBySlug(w http.ResponseWriter, r *http.Request) error {
	slug := chi.URLParam(r, "slug")
	room, err := h.queries.GetRoomBySlug(r.Context(), slug)
	if err != nil {
		return httpx.New(http.StatusNotFound, "not_found", "room not found", err)
	}
	return httpx.JSON(w, http.StatusOK, resdto.RoomRes{
		ID:        room.ID,
		Name:      room.Name,
		Slug:      room.Slug,
		CreatedAt: room.CreatedAt.Time,
	})
}

// Join handles joining a room.
func (h *RoomHandler) Join(w http.ResponseWriter, r *http.Request) error {
	claims, ok := auth.ClaimsFromCtx(r.Context())
	if !ok {
		return httpx.New(http.StatusUnauthorized, "unauthorized", "missing claims", nil)
	}

	roomID, err := strconv.ParseInt(chi.URLParam(r, "roomID"), 10, 64)
	if err != nil {
		return httpx.BadRequest("invalid_room_id", "invalid room id", err)
	}

	err = h.queries.JoinRoom(r.Context(), dbstore.JoinRoomParams{
		RoomID: roomID,
		UserID: claims.UserID,
	})
	if err != nil {
		return err
	}

	return httpx.JSON(w, http.StatusNoContent, nil)
}

// Leave handles leaving a room.
func (h *RoomHandler) Leave(w http.ResponseWriter, r *http.Request) error {
	claims, ok := auth.ClaimsFromCtx(r.Context())
	if !ok {
		return httpx.New(http.StatusUnauthorized, "unauthorized", "missing claims", nil)
	}

	roomID, err := strconv.ParseInt(chi.URLParam(r, "roomID"), 10, 64)
	if err != nil {
		return httpx.BadRequest("invalid_room_id", "invalid room id", err)
	}

	err = h.queries.LeaveRoom(r.Context(), dbstore.LeaveRoomParams{
		RoomID: roomID,
		UserID: claims.UserID,
	})
	if err != nil {
		return err
	}

	return httpx.JSON(w, http.StatusNoContent, nil)
}

// Messages handles fetching paginated message history for a room.
func (h *RoomHandler) Messages(w http.ResponseWriter, r *http.Request) error {
	roomID, err := strconv.ParseInt(chi.URLParam(r, "roomID"), 10, 64)
	if err != nil {
		return httpx.BadRequest("invalid_room_id", "invalid room id", err)
	}

	limit := int32(defaultMessageLimit)
	if l := r.URL.Query().Get("limit"); l != "" {
		n, err := strconv.ParseInt(l, 10, 32)
		if err == nil && n > 0 && n <= maxMessageLimit {
			limit = int32(n)
		}
	}

	msgs, err := h.queries.ListMessagesByRoom(r.Context(), dbstore.ListMessagesByRoomParams{
		RoomID: pgtype.Int8{Int64: roomID, Valid: true},
		Limit:  limit,
	})
	if err != nil {
		return err
	}

	res := make([]resdto.MessageRes, len(msgs))
	for i, m := range msgs {
		res[i] = resdto.MessageRes{
			ID:        m.ID,
			SenderID:  m.SenderID,
			Body:      m.Body,
			CreatedAt: m.CreatedAt.Time,
		}
		if m.RoomID.Valid {
			res[i].RoomID = &m.RoomID.Int64
		}
	}
	return httpx.JSON(w, http.StatusOK, res)
}

func slugify(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	s = strings.ReplaceAll(s, " ", "-")
	return s
}
