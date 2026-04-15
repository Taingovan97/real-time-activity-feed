// Package domain provides domain entities for the activity feed module.
package domain

import "time"

// FeedEvent represents a single activity feed event.
type FeedEvent struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Username  string    `json:"username"`
	EventType string    `json:"event_type"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}
