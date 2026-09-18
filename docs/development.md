# Development Guide

## Prerequisites

Install:

- Go 1.26.3 or a compatible Go toolchain.
- Node.js and npm.
- PostgreSQL.
- SQLC if you change SQL query definitions or schema-backed generated types.

## Environment

Create a `.env` file in the repository root. It is ignored by Git.

```env
DB_URL=postgres://user:password@localhost:5432/mylife?sslmode=disable
JWT_SECRET=replace-with-a-long-random-secret
```

`main.go` loads this file at startup. `DB_URL` is required to open PostgreSQL, and `JWT_SECRET` is required to sign and validate access tokens.

## Start the backend

From the repository root:

```sh
go run .
```

The API listens on `http://localhost:8080`.

## Start the frontend

In a second terminal:

```sh
cd frontend
npm install
npm run dev
```

The Vite development server normally listens on `http://localhost:5173`. The API CORS middleware is configured for that origin, and frontend requests currently target `http://localhost:8080` directly.

## Test the backend

From the repository root:

```sh
go test ./...
```

To run the API package tests with verbose output:

```sh
cd internal/api
go test -v
```

## Build and lint the frontend

```sh
cd frontend
npm run build
npm run lint
```

Use `npm run preview` after a build to serve the production bundle locally.

## SQLC workflow

SQL definitions live in `sql/queries` and schema files live in `sql/schema`. Generated Go code is written to `internal/database` according to `sqlc.yaml`.

After installing SQLC, regenerate the database package from the repository root with:

```sh
sqlc generate
```

Review generated changes together with the SQL change that caused them. Keep schema changes ordered and update queries and tests when a database contract changes.

## Common development checks

Before opening a pull request:

```sh
go test ./...
cd frontend && npm run lint && npm run build
```

The frontend and backend currently use hard-coded local development URLs. Update the API configuration and CORS policy together when running them on different hosts or ports.
