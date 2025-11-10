package store

import (
	"context"
	"database/sql"
	"time"
)

var (
	QueryTimeoutDuration = time.Second * 5
)

type Storage struct {
	UsersRepo Users
	PostsRepo Posts
}

type Posts interface {
	Create(context.Context, *Post) error
}
type Users interface {
	Create(context.Context, *User) error
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		PostsRepo: &PostStore{},
		UsersRepo: &UsersStore{},
	}
}
