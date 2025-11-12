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
	GetUserFeed(context.Context, int64) ([]PostsWithMetaData, error)
}

type Users interface {
	Create(context.Context, *User) error
	GetByID(context.Context, int64) (*User, error)
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
