# Ticket: DM REST endpoints

## Contexto

El cliente de direct messages ya está implementado. Para que funcione end-to-end necesita tres endpoints nuevos en el server. Las conversaciones ya se crean via WebSocket cuando se envía un DM, pero el cliente necesita poder listarlas, ver su historial, y buscar usuarios para iniciar nuevas.

## Trabajo requerido

### 1. SQL — `cmd/server/queries/`

- `messages.sql`: `ListMessagesByConversation` — sin JOIN a `users`, devuelve solo los campos de `messages`.
- `conversations.sql`: agregar query `ListConversationsByUser` que devuelva `(id, peer_id, peer_username)` — la otra persona en cada conversación del usuario autenticado.

Después de los cambios: `sqlc generate`.

### 2. Endpoints — todos protegidos con JWT

| Método | Path | Response |
|--------|------|----------|
| `GET` | `/api/v1/users?username=<q>` | `{id, username, created_at}` |
| `GET` | `/api/v1/conversations` | `[{id, peer_id, peer_username}]` |
| `GET` | `/api/v1/conversations/{id}/messages?limit=N` | `[{id, conversation_id, sender_id, body, created_at}]` |

## Criterios de aceptación

- Los tres endpoints responden correctamente con datos reales de la DB
- `ListMessagesByConversation` devuelve solo campos de `messages`, sin JOIN a `users`
- `GET /users` sin match devuelve 404
- `GET /conversations/{id}/messages` con un ID que no le pertenece al usuario autenticado devuelve 404 (no exponer conversaciones ajenas)
- El código sigue los patrones existentes: handler recibe `*dbstore.Queries`, usa `httpx.JSON`, errores via `httpx.New`

## Fuera de scope

- `POST /conversations` — se crea implícitamente via WS al enviar el primer DM
- Paginación con cursor — `limit` por query param es suficiente
