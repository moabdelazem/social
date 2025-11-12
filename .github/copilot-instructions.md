# Social API - Copilot Instructions

## Project Overview

Go REST API for social media platform using **chi router**, **PostgreSQL**, and **repository pattern** with interface-based storage.

## Architecture: 3-Layer Pattern

```
cmd/api/        → HTTP handlers, routing, JSON marshaling, middleware
internal/store/ → Repository interfaces + PostgreSQL implementations
internal/db/    → Connection pool setup, seeding utilities
```

**Critical Rule**: Database operations ONLY through `internal/store/` interfaces. Never call `db.Query()` directly in handlers.

**Store Organization**: Each domain has its own file (`posts.go`, `users.go`, `followers.go`, `comments.go`) with matching struct (`PostStore`, `UserStore`, etc.)

## Storage Layer Pattern

**All new repositories follow this pattern:**

1. **Define interface** in `internal/store/store.go`:

```go
type Posts interface {
    Create(context.Context, *Post) error
    GetByID(context.Context, int64) (*Post, error)
}
```

2. **Implement in separate file** (e.g., `posts.go`):

```go
type PostStore struct {
    db *sql.DB  // ← MUST be initialized or nil pointer panic
}
```

3. **Register in `NewStorage()`** - forgetting `db: db` causes runtime panics:

```go
func NewStorage(db *sql.DB) Storage {
    return Storage{
        PostsRepo: &PostStore{db: db},  // ← db parameter REQUIRED
    }
}
```

4. **Always use context timeouts** (already set to 5s via `QueryTimeoutDuration`):

```go
ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
defer cancel()
```

## Handler Development Pattern

**File**: `cmd/api/<resource>.go` (e.g., `posts.go`)

### 1. JSON Request/Response

- **Read**: `readJSON(w, r, &payload)` - enforces 1MB limit, rejects unknown fields
- **Write**: `app.jsonResponse(w, status, data)` - wraps in `{"data": ...}` envelope
- **Never** use `json.NewEncoder/Decoder` directly

### 2. Validation (go-playground/validator)

Use global `Validate` from `json.go`:

```go
type CreatePostPayload struct {
    Title string `json:"title" validate:"required,max=100"`
}
if err := Validate.Struct(payload); err != nil {
    app.badRequestResponse(w, r, err)  // ← Use error helpers
}
```

### 3. Error Handling

**Always use methods from `errors.go`** (never raw `writeJSONError`):

- `app.badRequestResponse(w, r, err)` - validation, parsing failures
- `app.notFoundResponse(w, r, err)` - `store.ErrorNotFound` from DB
- `app.internalServerError(w, r, err)` - database/server errors
- `app.conflictResponse(w, r, err)` - unique constraint violations

### 4. Status Codes

- `http.StatusCreated` (201) - POST creates
- `http.StatusOK` (200) - GET, PATCH
- `http.StatusNoContent` (204) - DELETE (no body)

## Routing (cmd/api/api.go)

**Chi router with nested routes**:

```go
r.Route("/v1", func(r chi.Router) {
    r.Route("/posts", func(r chi.Router) {
        r.Post("/", app.createPostHandler)

        r.Route("/{postID}", func(r chi.Router) {
            r.Use(app.postsContextMiddleware)  // ← Middleware loads resource
            r.Get("/", app.getPostHandler)
        })
    })
})
```

**URL params**: `chi.URLParam(r, "postID")` → `strconv.ParseInt()` → validate

**Custom error handlers** (configured in `mount()`):

- `r.NotFound()` - handles 404s
- `r.MethodNotAllowed()` - handles 405s

## Middleware Pattern

**File**: `cmd/api/middlewares.go` - Centralized middleware functions

**Context middleware** (see `postsContextMiddleware`, `usersContextMiddleware`):

1. Extract ID from URL params using `chi.URLParam(r, "postID")`
2. Parse with `strconv.ParseInt(idParam, 10, 64)`
3. Fetch resource from store
4. Handle `store.ErrorNotFound` with `errors.Is()` specifically
5. Store in context: `context.WithValue(ctx, "post", post)`
6. Retrieve later: `getPostFromCtx(r)` helper function

**Pattern**: Create matching getter functions for each context middleware (e.g., `getUserFromCtx`, `getPostFromCtx`)

## Database & Migrations

**Stack**: PostgreSQL 16 (Docker), golang-migrate, CloudBeaver (web UI)

**Start services**: `docker compose up -d`

- PostgreSQL: `localhost:5432` (credentials: `devuser/devpass`, database: `myapp_dev`)
- CloudBeaver: `localhost:8978` (web-based DB management, login: `admin/admin`)

### Migration Workflow (via Makefile)

```bash
make migrate-create NAME=add_followers   # Creates 000008_add_followers.{up,down}.sql
make migrate-up                           # Apply pending migrations
make migrate-down                         # Rollback last migration
make migrate-version                      # Check current version
make migrate-force VERSION=6              # Fix dirty state (when migration fails mid-run)
```

**Storage**: `cmd/migrate/migrations/` with sequential naming (`000001_`, `000002_`, etc.)

**Current migrations**: 8 total (users, posts, comments, followers, indexes)

**Never**: Run raw DDL via `db.Exec()` - always create migration files

## PostgreSQL Array Handling

Use `github.com/lib/pq` for array types (e.g., `tags TEXT[]`):

```go
// Insert
pq.Array(post.Tags)

// Scan
pq.Array(&post.Tags)
```

## Environment Configuration

**Package**: `internal/env/env.go` provides type-safe helpers:

```go
env.GetString("ADDR", ":6767")           // Fallback to :6767
env.GetInt("DB_MAX_OPEN_CONNS", 30)      // Parse int or fallback
```

**Loading**: `godotenv.Load()` in `main.go` reads `.env` (gitignored)

**Example .env**:

```
ADDR=:6767
DB_ADDR=postgres://devuser:devpass@localhost:5432/myapp_dev?sslmode=disable
```

## Development Commands

```bash
make watch   # Hot reload via Air (dev mode)
make build   # Compile to ./bin/main
make seed    # Populate test data (runs cmd/seed/main.go)
```

## Data Relationships & Eager Loading

**Pattern**: Fetch related data via separate repository calls:

```go
// In getPostHandler:
post := getPostFromCtx(r)
comments, _ := app.store.CommentRepo.GetByPostID(ctx, post.ID)
post.Comments = comments  // Attach to response
```

**Comments join users**: `CommentStore.GetByPostID()` uses `INNER JOIN users` to include user data.

**Feed aggregation**: `PostStore.GetUserFeed()` uses `LEFT JOIN` for comments count and `INNER JOIN followers` to filter posts from followed users.

## Error Handling Patterns

**Custom store errors** (defined in `store.go`):

- `ErrorNotFound` - Resource doesn't exist (return 404)
- `ErrorConflict` - Duplicate/constraint violation (return 409)
- `ErrorNotFollowing` - Unfollow operation on non-existent follow

**PostgreSQL error checking** (use `github.com/lib/pq`):

```go
if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
    return ErrorConflict  // Unique constraint violation
}
```

**Handler error switching**:

```go
switch {
case errors.Is(err, store.ErrorConflict):
    app.conflictResponse(w, r, err)
case errors.Is(err, store.ErrorNotFound):
    app.notFoundResponse(w, r, err)
default:
    app.internalServerError(w, r, err)
}
```

## Common Pitfalls

1. **Nil Pointer Panic in Store**: Check `NewStorage()` passes `db: db` to all struct fields
2. **405 Instead of 404**: Configure `r.NotFound()` and `r.MethodNotAllowed()` in `mount()`
3. **Hardcoded User IDs**: `POST /v1/posts` uses `UserID: 2` temporarily - run `make seed` first
4. **Array Scanning**: Always use `pq.Array()` for PostgreSQL array columns
5. **Context Leaks**: Always `defer cancel()` after `context.WithTimeout()`
6. **Duplicate Follows**: Use `ErrorConflict` detection for unique constraint violations

## Social Features Implementation

### Followers System

**Store**: `internal/store/followers.go` with dedicated `FollowerStore`

- `Follow(ctx, followerID, userID)` - Creates follow relationship
- `Unfollow(ctx, followerID, userID)` - Removes follow relationship
- Database handles self-follow prevention and cascading deletes

**Handlers**: `cmd/api/users.go`

- `PUT /v1/users/{userID}/follow` - Follow user (requires JSON body with `user_id`)
- `PUT /v1/users/{userID}/unfollow` - Unfollow user
- Both return `204 No Content` on success

### Feed System

**Store**: `PostStore.GetUserFeed(ctx, userID)` in `posts.go`

- Aggregates posts from followed users
- Includes comment counts via `LEFT JOIN`
- Returns `[]PostsWithMetaData` with `CommentsCount` field

**Handler**: `GET /v1/users/feed` (in `feed.go`, TODO: pagination)

## Key Files

- `cmd/api/api.go` - Router, middleware stack, custom error handlers
- `cmd/api/middlewares.go` - Context middleware (posts, users)
- `cmd/api/errors.go` - **USE THESE** for all error responses
- `cmd/api/json.go` - JSON helpers, validator init (`Validate`)
- `internal/store/store.go` - Repository interfaces, `NewStorage()` factory
- `internal/store/followers.go` - Follow/unfollow logic with conflict detection
- `internal/db/db.go` - Connection pool config (max conns, idle timeout)
- `Makefile` - All CLI commands (build, migrate, seed)
