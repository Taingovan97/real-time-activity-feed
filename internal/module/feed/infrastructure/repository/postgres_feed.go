// Package repository provides repository implementations for the feed module.
package repository

import (
	"context"
	"fmt"
	"strings"

	"real-time-activity-feed/internal/module/feed/application"
	"real-time-activity-feed/internal/module/feed/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresFeedRepository implements FeedRepository using PostgreSQL.
type PostgresFeedRepository struct {
	pool *pgxpool.Pool
}

// NewPostgresFeedRepository creates a PostgreSQL-backed feed repository.
func NewPostgresFeedRepository(pool *pgxpool.Pool) application.FeedRepository {
	return &PostgresFeedRepository{pool: pool}
}

// CreateEvent inserts a new feed event row.
func (r *PostgresFeedRepository) CreateEvent(
	ctx context.Context,
	entry domain.FeedEvent,
) (*domain.FeedEvent, error) {
	query := `
		INSERT INTO activity_feed_events (id, user_id, event_type, content)
		VALUES (uuid_generate_v4(), $1, $2, $3)
		RETURNING id, created_at
	`

	created := entry
	if err := r.pool.QueryRow(ctx, query, entry.UserID, entry.EventType, entry.Content).
		Scan(&created.ID, &created.CreatedAt); err != nil {
		return nil, fmt.Errorf("failed to create event: %w", err)
	}

	return &created, nil
}

// GetFeed returns feed items ordered newest first.
func (r *PostgresFeedRepository) GetFeed(
	ctx context.Context,
	eventType string,
	limit, offset int64,
) ([]domain.FeedEvent, int64, error) {
	queryBuilder := strings.Builder{}
	queryBuilder.WriteString(`
		SELECT
			e.id,
			e.user_id,
			u.username,
			e.event_type,
			e.content,
			e.created_at,
			COUNT(*) OVER() AS total
		FROM activity_feed_events e
		JOIN users u ON e.user_id = u.id
	`)

	args := []any{}
	if eventType != "" {
		queryBuilder.WriteString("\n\t\tWHERE e.event_type = $1")
		args = append(args, eventType)
	}

	queryBuilder.WriteString("\n\t\tORDER BY e.created_at DESC, e.id DESC")
	_, _ = fmt.Fprintf(&queryBuilder, "\n\t\tLIMIT $%d OFFSET $%d", len(args)+1, len(args)+2)
	args = append(args, limit, offset)

	rows, err := r.pool.Query(ctx, queryBuilder.String(), args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get feed: %w", err)
	}
	defer rows.Close()

	entries := make([]domain.FeedEvent, 0)
	var total int64
	for rows.Next() {
		var entry domain.FeedEvent
		if err := rows.Scan(
			&entry.ID,
			&entry.UserID,
			&entry.Username,
			&entry.EventType,
			&entry.Content,
			&entry.CreatedAt,
			&total,
		); err != nil {
			return nil, 0, fmt.Errorf("failed to scan feed item: %w", err)
		}
		entries = append(entries, entry)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating feed items: %w", err)
	}

	return entries, total, nil
}
