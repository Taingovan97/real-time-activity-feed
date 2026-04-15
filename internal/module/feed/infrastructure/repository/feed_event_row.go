// Package repository provides repository implementations for the feed module.
package repository

import "time"

// FeedEventRow represents a persisted activity feed event row.
type FeedEventRow struct {
	ID        string    `db:"id"`
	UserID    string    `db:"user_id"`
	EventType string    `db:"event_type"`
	Content   string    `db:"content"`
	CreatedAt time.Time `db:"created_at"`
}
