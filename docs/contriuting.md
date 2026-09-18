# Contributing

Thanks for improving MyLife. Keep changes focused, explain behavior changes, and include tests where the change affects the API or persistence layer.

## Where changes belong

- HTTP routes and request handling: `internal/api`.
- Authentication and token behavior: `internal/auth` and the authentication handlers in `internal/api`.
- SQL queries: `sql/queries`.
- Database schema changes: a new ordered file under `sql/schema`.
- Generated database code: `internal/database`; regenerate it from the SQL source instead of hand-editing it.
- Browser behavior and screens: `frontend/src`.
- Project-level documentation: `docs`.

## Workflow

1. Create a focused branch for the change.
2. Read the relevant handler, store interface, query, or React component before editing.
3. Make the smallest change that satisfies the requirement.
4. Add or update tests for changed behavior.
5. Run the backend and frontend checks from [development.md](development.md).
6. Review the diff for unrelated formatting or generated-file churn.
7. Open a pull request describing the behavior change and validation performed.

## API changes

When changing an endpoint, update all of the following when applicable:

- route registration in `internal/api/routes.go`;
- the handler and response model;
- frontend request code;
- handler tests;
- [api.md](api.md).

Preserve ownership checks: authenticated users must not be able to read, update, or delete another user's notes or profile.

## Database changes

Add schema changes as a new numbered SQL file. Update SQLC queries when needed, run `sqlc generate`, and test the affected handlers. Consider foreign keys, delete behavior, nullability, and existing data before changing a column.

## Testing expectations

Handler tests use store interfaces and fakes, which keeps most API tests independent of PostgreSQL. When adding a success-path test, stub every store method reached by the handler, not only the first database call. Run `go test ./...` before submitting the change.

## Commit guidance

Use a short imperative subject, such as:

```text
Document local development workflow
```

Keep commits reviewable and avoid mixing documentation, formatting, and unrelated refactors.
