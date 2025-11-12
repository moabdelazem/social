package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrorNotFound        = errors.New("record not found")
	ErrorConflict        = errors.New("resource already exists")
	ErrorNotFollowing    = errors.New("not following user")
	QueryTimeoutDuration = time.Second * 5
)

type Storage struct {
	UsersRepo    Users
	PostsRepo    Posts
	CommentRepo  Comments
	FollowerRepo Followers
}

type Posts interface {
	Create(context.Context, *Post) error
	GetByID(context.Context, int64) (*Post, error)
	Delete(context.Context, int64) error
	Update(context.Context, *Post) error
	GetUserFeed(context.Context, int64, PaginatedFeedQuery) ([]PostsWithMetaData, error)
}

type Users interface {
	Create(context.Context, *User) error
	GetByID(context.Context, int64) (*User, error)
	GetByEmail(context.Context, string) (*User, error)
	CreateAndInvite(context.Context, *User, string, time.Time) error
	Activate(context.Context, string) error
}

type Comments interface {
	GetByPostID(context.Context, int64) ([]Comment, error)
}

type Followers interface {
	Follow(ctx context.Context, followerID, userID int64) error
	Unfollow(ctx context.Context, followerID, userID int64) error
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		PostsRepo:    &PostStore{db: db},
		UsersRepo:    &UsersStore{db: db},
		CommentRepo:  &CommentStore{db: db},
		FollowerRepo: &FollowerStore{db: db},
	}
}

// withTx executes a function within a database transaction.
// If the function returns an error, the transaction is rolled back.
// If the function succeeds, the transaction is committed.
func withTx(db *sql.DB, ctx context.Context, fn func(*sql.Tx) error) error {
	// Start a new transaction
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// Defer rollback in case of panic or error
	// If Commit() is called successfully, Rollback() is a no-op
	defer tx.Rollback()

	// Execute the provided function with the transaction
	if err := fn(tx); err != nil {
		return err
	}

	// Commit the transaction if everything succeeded
	return tx.Commit()
}

// WithTx is the exported version of withTx for use outside the store package
func WithTx(db *sql.DB, ctx context.Context, fn func(*sql.Tx) error) error {
	return withTx(db, ctx, fn)
}
