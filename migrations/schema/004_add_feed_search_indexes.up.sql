CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE INDEX IF NOT EXISTS idx_activity_feed_events_event_type_created_id_desc
ON activity_feed_events (event_type, created_at DESC, id DESC);

CREATE INDEX IF NOT EXISTS idx_activity_feed_events_content_trgm
ON activity_feed_events USING GIN (content gin_trgm_ops);

CREATE INDEX IF NOT EXISTS idx_users_username_trgm
ON users USING GIN (username gin_trgm_ops);
