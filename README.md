# realtime-chat

Terminal-based real-time chat app. Rooms, direct messages, and WebSocket-powered messaging — all in the terminal.

## Stack

- **Server:** Go, [chi](https://github.com/go-chi/chi), PostgreSQL (pgx + sqlc + goose), WebSocket ([nhooyr.io/websocket](https://github.com/coder/websocket)), JWT auth
- **Client:** Go, [Bubble Tea](https://github.com/charmbracelet/bubbletea) TUI

## Getting started

### Server (Docker — recommended)

```bash
cp .env.example .env   # set JWT_SECRET at minimum
docker-compose up -d
```

This starts PostgreSQL and the server (with migrations applied automatically) on port 8080.

### Server (manual)

```bash
cp .env.example .env   # fill in DB_URL, PORT, JWT_SECRET
docker-compose up -d postgres
make migrate
make run-server
```

### Client

```bash
make run-client
```

## Features

- Register / login with JWT auth
- Create rooms, join/leave
- Real-time room messaging via WebSocket
- Direct messages (1-to-1)
- Message history (REST)
- Multiple themes: Catppuccin, Rose-Pine, Kanagawa

## Project layout

```
cmd/
  server/          # HTTP + WebSocket server
    internal/
      api/         # Router, middleware, handlers
      auth/        # JWT, bcrypt
      ws/          # WebSocket hub + client dispatch
      store/       # sqlc-generated DB layer
    migrations/    # goose migrations
    queries/       # SQL query files
  client/          # Bubble Tea TUI
    internal/
      api/         # HTTP client
      ws/          # WebSocket client
      ui/          # Screens (auth, rooms, chat, dm)
```

## WebSocket protocol

All messages use an envelope: `{"type": "<type>", "payload": {...}, "timestamp": "<RFC3339>"}`.

Client → server types: `room_message`, `direct_message`, `join_room`, `leave_room`, `user_typing`.
