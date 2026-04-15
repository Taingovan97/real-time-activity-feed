ALTER TABLE activity_feed_events
ALTER COLUMN created_at TYPE TIMESTAMPTZ
USING created_at AT TIME ZONE 'UTC';

ALTER TABLE activity_feed_events
ALTER COLUMN content TYPE TEXT;

CREATE INDEX IF NOT EXISTS idx_activity_feed_events_created_id_desc
ON activity_feed_events (created_at DESC, id DESC);
