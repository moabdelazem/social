package store

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"time"

	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int64    `json:"id"`
	Username  string   `json:"username"`
	Email     string   `json:"email"`
	Password  Password `json:"-"`
	IsActive  bool     `json:"is_active"`
	CreatedAt string   `json:"created_at"`
}

type Password struct {
	Text *string
	Hash []byte
}

func (p *Password) Set(text string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	p.Text = &text
	p.Hash = hash

	return nil
}

func (p *Password) ComparePassword(text string) error {
	return bcrypt.CompareHashAndPassword(p.Hash, []byte(text))
}

func (p *Password) Scan(value interface{}) error {
	if value == nil {
		p.Hash = nil
		return nil
	}

	switch v := value.(type) {
	case []byte:
		p.Hash = make([]byte, len(v))
		copy(p.Hash, v)
		return nil
	case string:
		p.Hash = []byte(v)
		return nil
	default:
		return nil
	}
}

type UsersStore struct {
	db *sql.DB
}

func (s *UsersStore) Create(ctx context.Context, user *User) error {
	query := `
		INSERT INTO users (username, email, password)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, is_active;
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(ctx, query, user.Username, user.Email, user.Password.Hash).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.IsActive,
	)

	if err != nil {
		// Check for unique constraint violation (duplicate username or email)
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return ErrorConflict
		}
		return err
	}

	return nil

}

func (s *UsersStore) GetByID(ctx context.Context, id int64) (*User, error) {
	query := `
		SELECT users.id, username, email, password, is_active, created_at FROM users
		WHERE id = $1;
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	user := &User{}
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.IsActive,
		&user.CreatedAt,
	)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrorNotFound
		default:
			return nil, err
		}
	}

	return user, nil
}

func (s *UsersStore) CreateAndInvite(ctx context.Context, user *User, token string, exp time.Time) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		// Create the user using the existing Create method logic
		if err := s.Create(ctx, user); err != nil {
			return err
		}

		// Create the invitation token
		if err := s.createUserInvitation(ctx, tx, token, exp, user.ID); err != nil {
			return err
		}
		return nil
	})
}

func (s *UsersStore) createUserInvitation(ctx context.Context, tx *sql.Tx, token string, exp time.Time, userID int64) error {
	query := `
		INSERT INTO user_invitations (token, user_id, expiry)
		VALUES ($1, $2, $3)	
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, token, userID, exp)
	if err != nil {
		return err
	}

	return nil
}

func (s *UsersStore) Activate(ctx context.Context, token string) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		user, err := s.getUserFromInvitation(ctx, tx, token)
		if err != nil {
			return err
		}

		user.IsActive = true
		if err := s.update(ctx, tx, user); err != nil {
			return err
		}

		if err := s.deleteUserInvitations(ctx, tx, user.ID); err != nil {
			return err
		}

		return nil
	})
}

func (s *UsersStore) getUserFromInvitation(ctx context.Context, tx *sql.Tx, token string) (*User, error) {
	query := `
		SELECT u.id, u.username, u.email, u.created_at, u.is_active
		FROM users u
		JOIN user_invitations ui ON u.id = ui.user_id
		WHERE ui.token = $1 AND ui.expiry > $2
	`

	hash := sha256.Sum256([]byte(token))
	hashToken := hex.EncodeToString(hash[:])

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	user := &User{}
	err := tx.QueryRowContext(ctx, query, hashToken, time.Now()).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.CreatedAt,
		&user.IsActive,
	)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrorNotFound
		default:
			return nil, err
		}
	}

	return user, nil
}

func (s *UsersStore) update(ctx context.Context, tx *sql.Tx, user *User) error {
	query := `UPDATE users SET username = $1, email = $2, is_active = $3 WHERE id = $4`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, user.Username, user.Email, user.IsActive, user.ID)
	if err != nil {
		return err
	}

	return nil
}

func (s *UsersStore) deleteUserInvitations(ctx context.Context, tx *sql.Tx, userID int64) error {
	query := `DELETE FROM user_invitations WHERE user_id = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, userID)
	if err != nil {
		return err
	}

	return nil
}

func (s *UsersStore) delete(ctx context.Context, tx *sql.Tx, id int64) error {
	query := `DELETE FROM users WHERE id = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *UsersStore) GetByEmail(ctx context.Context, email string) (*User, error) {
	query := `
		SELECT id, username, email, password, created_at FROM users
		WHERE email = $1 AND is_active = true
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	user := &User{}
	err := s.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password.Hash,
		&user.CreatedAt,
	)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrorNotFound
		default:
			return nil, err
		}
	}

	return user, nil
}
