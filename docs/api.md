# API Reference

The API listens on `http://localhost:8080` by default. JSON endpoints return JSON bodies. Protected endpoints require:

```http
Authorization: Bearer <access-token>
```

Access tokens are JWTs signed with `JWT_SECRET` and expire after one hour. The frontend stores the access token in browser local storage under `accessToken`.

## Authentication

### Register

`POST /Dashboard/RegisterUser`

Request:

```json
{
  "user_name": "ada",
  "email": "ada@example.com",
  "password": "correct horse"
}
```

Passwords must be at least six characters. Success returns `201 Created`:

```json
{
  "id": "uuid",
  "email": "ada@example.com",
  "user_name": "ada"
}
```

### Login

`POST /Dashboard/LoginUser`

Request fields are the same as registration. Success returns `200 OK`:

```json
{
  "id": "uuid",
  "user_name": "ada",
  "email": "ada@example.com",
  "created_at": "timestamp",
  "access_token": "jwt",
  "refresh_token": "opaque-token"
}
```

### Refresh an access token

`POST /api/refresh`

Send the refresh token in the bearer header. Success returns `200 OK`:

```json
{
  "token": "jwt"
}
```

### Revoke a refresh token

`POST /api/revoke`

Send the refresh token in the bearer header. Success returns `204 No Content`.

## Profile and notes

The `username` path value is used to locate the profile, but the authenticated JWT subject is the identity used for authorization. A user can only access their own profile and notes.

### Get profile

`GET /Profile/{username}`

Returns `200 OK` with the profile, streak counters, note counts, and up to three recent notes.

### Get all notes

`GET /Profile/{username}/GetNotes`

Returns `200 OK` with an array of notes:

```json
[
  {
    "note_id": "uuid",
    "user_id": "uuid",
    "daily_note": "Wrote project documentation",
    "created_at": "timestamp",
    "last_updated": "timestamp"
  }
]
```

### Create a note

`POST /Profile/{username}/CreateNote`

Request:

```json
{
  "daily_note": "Wrote project documentation"
}
```

The handler creates the note, updates streak and note counters, and returns the created note with `200 OK`.

### Read a note

`GET /Profile/{username}/ReadNote/{note_id}`

Returns the note with `200 OK`. Common errors are `400 Bad Request` for an invalid UUID, `401 Unauthorized` for a missing or invalid token, `403 Forbidden` for another user's note, and `404 Not Found` when the note does not exist.

### Update a note

`PUT /Profile/{username}/UpdateNote/{note_id}`

Request:

```json
{
  "daily_note": "Updated text"
}
```

Returns the updated note with `200 OK`.

### Delete a note

`DELETE /Profile/{username}/DeleteNote/{note_id}`

Returns `200 OK` with a success message when deletion succeeds.

## Error responses

Errors use this response shape:

```json
{
  "error: ": "Not Authorized"
}
```

The trailing space and colon are part of the current response key for compatibility with the implementation.

## CORS

The API currently allows requests from `http://localhost:5173` and accepts `Content-Type` and `Authorization` headers. Update `internal/api/middleware.go` when changing the frontend origin.
