package application

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"real-time-activity-feed/internal/module/feed/domain"
	"real-time-activity-feed/internal/module/feed/infrastructure/mocks"
	"real-time-activity-feed/internal/shared/logger"
)

func TestFeedUseCase_GetFeed_FromCache(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cache := mocks.NewMockFeedCacheRepository(ctrl)
	persistence := mocks.NewMockFeedRepository(ctrl)
	broadcast := mocks.NewMockBroadcastService(ctrl)
	users := mocks.NewMockUserRepository(ctrl)

	cache.EXPECT().GetFeed(ctx, int64(10), int64(0)).Return([]domain.FeedEvent{{ID: "evt-1"}}, int64(1), nil)

	uc := NewFeedUseCase(cache, persistence, users, broadcast, logger.New("info", false))
	entries, total, err := uc.GetFeed(ctx, 10, 0)

	require.NoError(t, err)
	require.Len(t, entries, 1)
	require.Equal(t, int64(1), total)
}

func TestFeedUseCase_GetFeed_FallbackToPersistence(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cache := mocks.NewMockFeedCacheRepository(ctrl)
	persistence := mocks.NewMockFeedRepository(ctrl)
	broadcast := mocks.NewMockBroadcastService(ctrl)
	users := mocks.NewMockUserRepository(ctrl)

	cache.EXPECT().GetFeed(ctx, int64(10), int64(0)).Return(nil, int64(0), errors.New("redis down"))
	persistence.EXPECT().GetFeed(ctx, int64(10), int64(0)).Return([]domain.FeedEvent{{ID: "evt-1"}}, int64(1), nil)

	uc := NewFeedUseCase(cache, persistence, users, broadcast, logger.New("info", false))
	entries, total, err := uc.GetFeed(ctx, 10, 0)

	require.NoError(t, err)
	require.Len(t, entries, 1)
	require.Equal(t, int64(1), total)
}
