# API Draft

## Health

- `GET /health`

Response:

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "status": "ok"
  }
}
```

## Auth

- `POST /api/v1/auth/register`
- `POST /api/v1/auth/login`

Register request:

```json
{
  "username": "alice",
  "password": "123456"
}
```

Login response:

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "token": "demo-token-1",
    "type": "Bearer"
  }
}
```

## Tasks

- `POST /api/v1/tasks`
- `GET /api/v1/tasks`
- `GET /api/v1/tasks/:id`
- `PUT /api/v1/tasks/:id`
- `DELETE /api/v1/tasks/:id`

Create request:

```json
{
  "title": "learn Go API",
  "content": "build the first simple version",
  "status": "todo"
}
```

Protected task write APIs currently accept demo tokens:

```text
Authorization: Bearer demo-token-1
```
