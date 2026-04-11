## Ticket activo — DM endpoints

- [ ] `sqlc generate`
- [ ] `GET /api/v1/conversations`
- [ ] `GET /api/v1/conversations/{id}/messages`
- [ ] `GET /api/v1/users?username=xxx`

## Refactors de rooms

- [ ] Migrar `Create`, `List`, `Join`, `Leave`, `Messages` al service — sacar `h.queries` del handler
- [ ] `room.Service` devolver tipo de dominio, no `response.RoomRes`
- [ ] `NewService` no asigna el logger
- [ ] `pgx.ErrNoRows` wrapearlo en el service, no en el handler

## Mejoras generales

- [ ] Paginación en endpoints de listado (`/rooms` y otros)
- [ ] Revisión general de clean architecture
- [ ] `go-playground/validator` para reemplazar validación manual

## Features futuras

- [ ] Notification inbox (requiere refactor del ciclo de vida del WS primero)
- [ ] Paquete `pkg/dto/` compartido entre server y client
