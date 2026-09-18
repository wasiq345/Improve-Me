# MyLife Documentation

MyLife is a daily-notes application with a Go HTTP API, a React/Vite frontend, and PostgreSQL persistence.

## Documents

- [API reference](api.md): HTTP routes, authentication, request bodies, and response shapes.
- [Architecture](architecture.md): runtime components, data flow, and persistence model.
- [Development guide](development.md): prerequisites, environment variables, local commands, and testing.
- [Contributor guide](contriuting.md): contribution workflow, code organization, and change checklist.

## Project layout

```text
.
|-- main.go                 Go server entry point
|-- internal/
|   |-- api/                HTTP handlers, routes, middleware, and tests
|   |-- auth/               JWT, refresh-token, and password helpers
|   `-- database/           SQLC-generated database access
|-- sql/
|   |-- queries/            SQLC query definitions
|   `-- schema/             PostgreSQL schema changes
|-- frontend/               React/Vite client
`-- docs/                   Project documentation
```
