package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrorNotFound        = errors.New("record not found")
	QueryTimeoutDuration = time.Second * 5
)

type Storage struct {
	UsersRepo Users
	PostsRepo Posts
}

type Posts interface {
	Create(context.Context, *Post) error
	GetByID(context.Context, int64) (*Post, error)
	Delete(context.Context, int64) error
}
type Users interface {
	Create(context.Context, *User) error
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		PostsRepo: &PostStore{db: db},
		UsersRepo: &UsersStore{db: db},
	}
}
