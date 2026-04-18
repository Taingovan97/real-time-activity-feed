// Package cache provides feed cache implementations.
package cache

import (
	"context"

	"real-time-activity-feed/internal/module/feed/domain"
)

// NoopFeedCache disables feed caching while satisfying the FeedCache contract.
type NoopFeedCache struct{}

// NewNoopFeedCache creates a cache implementation that always bypasses storage.
func NewNoopFeedCache() *NoopFeedCache {
	return &NoopFeedCache{}
}

// GetFeed always reports a cache miss.
func (c *NoopFeedCache) GetFeed(_ context.Context, _ string, _ int64, _ int64) ([]domain.FeedEvent, int64, bool, error) {
	return nil, 0, false, nil
}

// SetFeed ignores writes and returns success.
func (c *NoopFeedCache) SetFeed(_ context.Context, _ string, _ int64, _ int64, _ []domain.FeedEvent, _ int64) error {
	return nil
}

// InvalidateAfterPublish ignores invalidation requests and returns success.
func (c *NoopFeedCache) InvalidateAfterPublish(_ context.Context, _ string) error {
	return nil
}
