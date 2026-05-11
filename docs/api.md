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

## Users

- `PUT /api/v1/users/me`

Update current user request:

```json
{
  "username": "alice_new",
  "password": "abcdef"
}
```

The request requires the demo token returned by login. `username` and
`password` are both optional, but at least one field must be provided.

## Tasks

- `POST /api/v1/tasks`
- `GET /api/v1/tasks`
- `GET /api/v1/tasks/:id`
- `PUT /api/v1/tasks/:id`
- `DELETE /api/v1/tasks/:id`

All task APIs require the demo token returned by login. A user can only list,
read, update, and delete their own tasks.

Create request:

```json
{
  "title": "learn Go API",
  "content": "build the first simple version",
  "status": "todo"
}
```

Task APIs currently accept demo tokens:

```text
Authorization: Bearer demo-token-1
```
