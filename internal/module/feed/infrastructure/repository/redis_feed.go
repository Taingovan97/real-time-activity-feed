// Package repository provides repository implementations for the feed module.
package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"real-time-activity-feed/internal/module/feed/application"
	"real-time-activity-feed/internal/module/feed/domain"

	"github.com/redis/go-redis/v9"
)

// RedisFeedRepository implements FeedCacheRepository using a Redis list of recent feed items.
type RedisFeedRepository struct {
	client *redis.Client
}

// NewRedisFeedRepository creates a Redis-backed recent feed cache.
func NewRedisFeedRepository(client *redis.Client) application.FeedCacheRepository {
	return &RedisFeedRepository{client: client}
}

// AddEvent prepends a feed item to the recent-events cache.
func (r *RedisFeedRepository) AddEvent(ctx context.Context, entry domain.FeedEvent) error {
	payload, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("failed to marshal feed item: %w", err)
	}

	pipe := r.client.TxPipeline()
	pipe.LPush(ctx, domain.RedisFeedKey, payload)
	pipe.LTrim(ctx, domain.RedisFeedKey, 0, domain.MaxFeedItems-1)
	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("failed to cache feed item: %w", err)
	}

	return nil
}

// GetFeed returns cached feed items ordered newest first.
func (r *RedisFeedRepository) GetFeed(
	ctx context.Context,
	limit, offset int64,
) ([]domain.FeedEvent, int64, error) {
	total, err := r.client.LLen(ctx, domain.RedisFeedKey).Result()
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get feed size: %w", err)
	}

	if total == 0 {
		return []domain.FeedEvent{}, 0, nil
	}

	start := offset
	stop := offset + limit - 1
	results, err := r.client.LRange(ctx, domain.RedisFeedKey, start, stop).Result()
	if err != nil {
		return nil, 0, fmt.Errorf("failed to read feed cache: %w", err)
	}

	entries := make([]domain.FeedEvent, 0, len(results))
	for _, item := range results {
		var entry domain.FeedEvent
		if err := json.Unmarshal([]byte(item), &entry); err != nil {
			continue
		}
		entries = append(entries, entry)
	}

	return entries, total, nil
}
