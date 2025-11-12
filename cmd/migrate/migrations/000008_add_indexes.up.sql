-- Enable pg_trgm extension for text search
CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- Users table indexes
CREATE INDEX IF NOT EXISTS idx_users_username ON users (username);

CREATE INDEX IF NOT EXISTS idx_users_email ON users (email);

-- Posts table indexes
CREATE INDEX IF NOT EXISTS idx_posts_user_id ON posts (user_id);

CREATE INDEX IF NOT EXISTS idx_posts_created_at ON posts (created_at DESC);

CREATE INDEX IF NOT EXISTS idx_posts_tags ON posts USING GIN (tags);

CREATE INDEX IF NOT EXISTS idx_posts_title ON posts USING GIN (title gin_trgm_ops);

-- Comments table indexes
CREATE INDEX IF NOT EXISTS idx_comments_post_id ON comments (post_id);

CREATE INDEX IF NOT EXISTS idx_comments_user_id ON comments (user_id);

CREATE INDEX IF NOT EXISTS idx_comments_created_at ON comments (created_at DESC);

-- Followers table indexes (composite primary key already indexed)
CREATE INDEX IF NOT EXISTS idx_followers_follower_id ON followers (follower_id);

CREATE INDEX IF NOT EXISTS idx_followers_created_at ON followers (created_at DESC);