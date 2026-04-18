package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"real-time-activity-feed/internal/module/feed/domain"
)

type cachedFeedPage struct {
	Entries []domain.FeedEvent `json:"entries"`
	Total   int64              `json:"total"`
}

// RedisFeedCache stores hot feed pages in Redis.
type RedisFeedCache struct {
	client *redis.Client
	ttl    time.Duration
}

func NewRedisFeedCache(client *redis.Client, ttl time.Duration) *RedisFeedCache {
	return &RedisFeedCache{
		client: client,
		ttl:    ttl,
	}
}

func buildFeedPageKey(eventType string, limit, offset int64) string {
	if eventType == "" {
		return fmt.Sprintf("feed:v1:page:all:limit=%d:offset=%d", limit, offset)
	}

	return fmt.Sprintf("feed:v1:page:type=%s:limit=%d:offset=%d", eventType, limit, offset)
}

func marshalCachedFeedPage(page cachedFeedPage) (string, error) {
	data, err := json.Marshal(page)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func unmarshalCachedFeedPage(raw string) (cachedFeedPage, error) {
	var page cachedFeedPage
	err := json.Unmarshal([]byte(raw), &page)
	return page, err
}

func (c *RedisFeedCache) GetFeed(ctx context.Context, eventType string, limit, offset int64) ([]domain.FeedEvent, int64, bool, error) {
	raw, err := c.client.Get(ctx, buildFeedPageKey(eventType, limit, offset)).Result()
	if err == redis.Nil {
		return nil, 0, false, nil
	}
	if err != nil {
		return nil, 0, false, err
	}

	page, err := unmarshalCachedFeedPage(raw)
	if err != nil {
		return nil, 0, false, err
	}

	return page.Entries, page.Total, true, nil
}

func (c *RedisFeedCache) SetFeed(ctx context.Context, eventType string, limit, offset int64, entries []domain.FeedEvent, total int64) error {
	payload, err := marshalCachedFeedPage(cachedFeedPage{
		Entries: entries,
		Total:   total,
	})
	if err != nil {
		return err
	}

	return c.client.Set(ctx, buildFeedPageKey(eventType, limit, offset), payload, c.ttl).Err()
}

func (c *RedisFeedCache) InvalidateAfterPublish(ctx context.Context, eventType string) error {
	keys := []string{
		buildFeedPageKey("", 10, 0),
		buildFeedPageKey(eventType, 10, 0),
	}

	return c.client.Del(ctx, keys...).Err()
}
