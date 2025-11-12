package store

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"time"

	"github.com/lib/pq"
)

type Post struct {
	ID        int64     `json:"id"`
	Content   string    `json:"content"`
	Title     string    `json:"title"`
	UserID    int64     `json:"user_id"`
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Comments  []Comment `json:"comments"`
	Version   int       `json:"version"`
}

type PostsWithMetaData struct {
	Post
	CommentsCount int `json:"comments_count"`
}

type PostStore struct {
	db *sql.DB
}

func (s *PostStore) GetUserFeed(ctx context.Context, userId int64, fq PaginatedFeedQuery) ([]PostsWithMetaData, error) {
	// Build dynamic query with filters
	query := `
		SELECT 
			p.id, p.user_id, p.title, p.content, p.created_at, p.updated_at, p.tags, p.version,
			COUNT(c.id) AS comments_count
		FROM posts p
		LEFT JOIN comments c ON c.post_id = p.id
		INNER JOIN followers f ON f.user_id = p.user_id
		WHERE f.follower_id = $1`

	// Dynamic query params
	args := []interface{}{userId}
	paramIndex := 2

	// Add search filter (search in title and content)
	if fq.Search != "" {
		query += ` AND (p.title ILIKE $` + strconv.Itoa(paramIndex) + ` OR p.content ILIKE $` + strconv.Itoa(paramIndex) + `)`
		args = append(args, "%"+fq.Search+"%")
		paramIndex++
	}

	// Add tags filter (posts must contain all specified tags)
	if len(fq.Tags) > 0 {
		query += ` AND p.tags @> $` + strconv.Itoa(paramIndex)
		args = append(args, pq.Array(fq.Tags))
		paramIndex++
	}

	// Add since filter (posts created after this date)
	if fq.Since != "" {
		query += ` AND p.created_at >= $` + strconv.Itoa(paramIndex)
		args = append(args, fq.Since)
		paramIndex++
	}

	// Add until filter (posts created before this date)
	if fq.Until != "" {
		query += ` AND p.created_at <= $` + strconv.Itoa(paramIndex)
		args = append(args, fq.Until)
		paramIndex++
	}

	// Add GROUP BY, ORDER BY, LIMIT, and OFFSET
	query += `
		GROUP BY p.id
		ORDER BY p.created_at ` + fq.Sort + `
		LIMIT $` + strconv.Itoa(paramIndex) + ` OFFSET $` + strconv.Itoa(paramIndex+1)

	args = append(args, fq.Limit, fq.Offset)

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	feed := make([]PostsWithMetaData, 0)
	for rows.Next() {
		var p PostsWithMetaData
		err := rows.Scan(
			&p.ID,
			&p.UserID,
			&p.Title,
			&p.Content,
			&p.CreatedAt,
			&p.UpdatedAt,
			pq.Array(&p.Tags),
			&p.Version,
			&p.CommentsCount,
		)
		if err != nil {
			return nil, err
		}
		feed = append(feed, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return feed, nil
}

func (s *PostStore) Create(ctx context.Context, post *Post) error {
	query := `
	INSERT INTO posts (content, title, user_id, tags)
	VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at, version
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		post.Content,
		post.Title,
		post.UserID,
		pq.Array(post.Tags),
	).Scan(
		&post.ID,
		&post.CreatedAt,
		&post.UpdatedAt,
		&post.Version,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *PostStore) GetByID(ctx context.Context, id int64) (*Post, error) {
	query := `
		SELECT id, user_id, title, content, created_at, updated_at, tags, version
	 	FROM posts
		WHERE id = $1
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var post Post
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&post.ID,
		&post.UserID,
		&post.Title,
		&post.Content,
		&post.CreatedAt,
		&post.UpdatedAt,
		pq.Array(&post.Tags),
		&post.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrorNotFound
		default:
			return nil, err
		}
	}

	return &post, nil
}

func (s *PostStore) Delete(ctx context.Context, postID int64) error {
	query := `
		DELETE FROM posts WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	res, err := s.db.ExecContext(ctx, query, postID)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrorNotFound
	}

	return nil
}

func (s *PostStore) Update(ctx context.Context, post *Post) error {
	query := `
		UPDATE posts
		SET title = $1, content = $2, version = version + 1, updated_at = NOW()
		WHERE id = $3 AND version = $4
		RETURNING version, updated_at
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		post.Title,
		post.Content,
		post.ID,
		post.Version,
	).Scan(
		&post.Version,
		&post.UpdatedAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrorNotFound
		default:
			return err
		}
	}

	return nil
}
