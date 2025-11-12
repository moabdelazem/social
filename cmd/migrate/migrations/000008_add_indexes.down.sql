-- Drop followers table indexes
DROP INDEX IF EXISTS idx_followers_created_at;

DROP INDEX IF EXISTS idx_followers_follower_id;

-- Drop comments table indexes
DROP INDEX IF EXISTS idx_comments_created_at;

DROP INDEX IF EXISTS idx_comments_user_id;

DROP INDEX IF EXISTS idx_comments_post_id;

-- Drop posts table indexes
DROP INDEX IF EXISTS idx_posts_title_trgm;

DROP INDEX IF EXISTS idx_posts_tags;

DROP INDEX IF EXISTS idx_posts_created_at;

DROP INDEX IF EXISTS idx_posts_user_id;

-- Drop users table indexes
DROP INDEX IF EXISTS idx_users_email;

DROP INDEX IF EXISTS idx_users_username;

-- Drop extension (optional - keep if other migrations might use it)
-- DROP EXTENSION IF EXISTS pg_trgm;