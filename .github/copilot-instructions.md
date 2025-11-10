# Social API - Copilot Instructions

## Project Overview

Go REST API for a social media platform using chi router, PostgreSQL, and a repository pattern with interface-based storage layer.

## Architecture Pattern: Repository + Interface Storage

**Critical**: All database operations go through interface-based repositories in `internal/store/`:

- Define interface in `store.go` (e.g., `Posts`, `Users`, `Comments`)
- Implement in separate files (`posts.go`, `users.go`, `comments.go`)
- **Must pass `db *sql.DB` when initializing stores** in `NewStorage()` - forgetting this causes nil pointer panics
- All store methods use `context.WithTimeout(ctx, QueryTimeoutDuration)` for database operations

```go
// Example: Always initialize stores with db connection
func NewStorage(db *sql.DB) Storage {
    return Storage{
        PostsRepo: &PostStore{db: db},  // ← db parameter is REQUIRED
        UsersRepo: &UsersStore{db: db},
    }
}
```

## Handler Pattern (cmd/api/)

1. **Error Handling**: Use centralized error methods from `errors.go`, never direct `writeJSONError()`:

   - `app.badRequestResponse(w, r, err)` - validation/parsing errors
   - `app.notFoundResponse(w, r, err)` - missing resources
   - `app.internalServerError(w, r, err)` - database/server errors
   - `app.conflictResponse(w, r, err)` - duplicates/conflicts

2. **Validation**: Use `go-playground/validator` via global `Validate` (initialized in `json.go`):

```go
type CreatePayload struct {
    Title string `json:"title" validate:"required,max=100"`
}
// Then: if err := Validate.Struct(payload); err != nil
```

3. **JSON Operations**: Use helper functions from `json.go`:

   - `readJSON(w, r, &payload)` - max 1MB, disallows unknown fields
   - `writeJSON(w, status, data)` - auto sets Content-Type
   - Never use direct `json.Encoder/Decoder` in handlers

4. **Response Status Codes**:
   - POST create: `http.StatusCreated` (201)
   - GET: `http.StatusOK` (200)
   - DELETE: `http.StatusNoContent` (204)

## Database & Migrations

**Setup**: PostgreSQL via Docker Compose (`docker compose up -d`)

- Connection: `postgres://devuser:devpass@localhost:5432/myapp_dev?sslmode=disable`
- Use `golang-migrate` for schema changes (NOT manual SQL in `db.Exec`)

**Migration Commands** (via Makefile):

```bash
make migrate-create NAME=add_feature    # Creates up/down pair
make migrate-up                          # Apply migrations
make migrate-down                        # Rollback one migration
make migrate-version                     # Check current version
make migrate-force VERSION=X             # Fix dirty migrations
```

**Migration Pattern**: Stored in `cmd/migrate/migrations/`, sequential naming `NNNNNN_description.{up,down}.sql`

## Development Workflow

**Run dev server**: `make watch` (uses Air for hot reload)
**Build binary**: `make build` → `./bin/main`
**Config**: Environment vars in `.env` (loaded via `godotenv` in `main.go`)

## Common Pitfalls

1. **Nil DB Panic**: If you see "nil pointer dereference" in store methods, check `NewStorage()` passes `db: db` to all stores
2. **Foreign Keys**: Posts reference `user_id` → users table. Create test users before posts
3. **Array Types**: PostgreSQL arrays use `pq.Array()` for tags: `pq.Array(post.Tags)` when inserting, `pq.Array(&post.Tags)` when scanning
4. **Context Timeouts**: Always create timeout contexts in store methods before DB calls (see `QueryTimeoutDuration = 5s`)
5. **Chi URL Params**: Use `chi.URLParam(r, "postID")` for path params, then `strconv.ParseInt()`

## Testing Data Flow

Posts with comments example:

- `POST /v1/posts` creates post (requires `user_id` in body, currently hardcoded to 2)
- `GET /v1/posts/{id}` returns post WITH comments array (joined query in `CommentStore.GetByPostID`)
- Comments fetched via separate repository method after post retrieval

## Key Files Reference

- `cmd/api/main.go` - Entry point, config, DB initialization
- `cmd/api/api.go` - Router setup (chi), middleware configuration
- `cmd/api/errors.go` - Centralized error handlers (USE THESE, not direct JSON writes)
- `internal/store/store.go` - Storage interfaces and factory
- `internal/db/db.go` - DB connection pool setup
- `.env` - Database connection strings (gitignored)
