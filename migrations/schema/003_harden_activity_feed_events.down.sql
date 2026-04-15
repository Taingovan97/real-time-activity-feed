DROP INDEX IF EXISTS idx_activity_feed_events_created_id_desc;

ALTER TABLE activity_feed_events
ALTER COLUMN created_at TYPE TIMESTAMP
USING created_at AT TIME ZONE 'UTC';

ALTER TABLE activity_feed_events
ALTER COLUMN content TYPE VARCHAR(255);
