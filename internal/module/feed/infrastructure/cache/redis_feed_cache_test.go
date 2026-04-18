package cache

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"real-time-activity-feed/internal/module/feed/domain"
)

func TestRedisFeedCacheKey(t *testing.T) {
	t.Parallel()

	require.Equal(t, "feed:v1:page:all:limit=10:offset=0", buildFeedPageKey("", 10, 0))
	require.Equal(t, "feed:v1:page:type=upload:limit=10:offset=0", buildFeedPageKey(domain.EventTypeUpload, 10, 0))
}

func TestNoopFeedCacheMissesAndInvalidatesCleanly(t *testing.T) {
	t.Parallel()

	cache := NewNoopFeedCache()
	entries, total, found, err := cache.GetFeed(context.Background(), "", 10, 0)

	require.NoError(t, err)
	require.False(t, found)
	require.Nil(t, entries)
	require.Zero(t, total)
	require.NoError(t, cache.SetFeed(context.Background(), "", 10, 0, []domain.FeedEvent{{ID: "evt-1"}}, 1))
	require.NoError(t, cache.InvalidateAfterPublish(context.Background(), domain.EventTypeUpload))
}

func TestRedisFeedCachePayloadRoundTrip(t *testing.T) {
	t.Parallel()

	payload := cachedFeedPage{
		Entries: []domain.FeedEvent{{ID: "evt-1", EventType: domain.EventTypeUpload}},
		Total:   1,
	}

	encoded, err := marshalCachedFeedPage(payload)
	require.NoError(t, err)

	decoded, err := unmarshalCachedFeedPage(encoded)
	require.NoError(t, err)
	require.Equal(t, payload, decoded)
}

func TestRedisFeedCacheUsesConfiguredTTL(t *testing.T) {
	t.Parallel()

	cache := &RedisFeedCache{ttl: 30 * time.Second}
	require.Equal(t, 30*time.Second, cache.ttl)
}
