package cache

import (
	"context"

	"real-time-activity-feed/internal/module/feed/domain"
)

// NoopFeedCache disables feed caching while satisfying the FeedCache contract.
type NoopFeedCache struct{}

func NewNoopFeedCache() *NoopFeedCache {
	return &NoopFeedCache{}
}

func (c *NoopFeedCache) GetFeed(_ context.Context, _ string, _ int64, _ int64) ([]domain.FeedEvent, int64, bool, error) {
	return nil, 0, false, nil
}

func (c *NoopFeedCache) SetFeed(_ context.Context, _ string, _ int64, _ int64, _ []domain.FeedEvent, _ int64) error {
	return nil
}

func (c *NoopFeedCache) InvalidateAfterPublish(_ context.Context, _ string) error {
	return nil
}
