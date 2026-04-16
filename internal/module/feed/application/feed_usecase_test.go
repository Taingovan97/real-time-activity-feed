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

	persistence := mocks.NewMockFeedRepository(ctrl)
	broadcast := mocks.NewMockBroadcastService(ctrl)
	users := mocks.NewMockUserRepository(ctrl)

	persisted := []domain.FeedEvent{{ID: "evt-1"}}
	persistence.EXPECT().GetFeed(ctx, "", int64(10), int64(0)).Return(persisted, int64(1), nil)

	uc := NewFeedUseCase(persistence, users, broadcast, logger.New("info", false))
	entries, total, err := uc.GetFeed(ctx, "", 10, 0)

	require.NoError(t, err)
	require.Equal(t, persisted, entries)
	require.Equal(t, int64(1), total)
}

func TestFeedUseCase_GetFeed_WhenPersistenceFails(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	persistence := mocks.NewMockFeedRepository(ctrl)
	broadcast := mocks.NewMockBroadcastService(ctrl)
	users := mocks.NewMockUserRepository(ctrl)

	persistence.EXPECT().GetFeed(ctx, "", int64(10), int64(0)).Return(nil, int64(0), errors.New("db down"))

	uc := NewFeedUseCase(persistence, users, broadcast, logger.New("info", false))
	entries, total, err := uc.GetFeed(ctx, "", 10, 0)

	require.Error(t, err)
	require.Nil(t, entries)
	require.Zero(t, total)
}

func TestFeedUseCase_GetFeed_WithEventTypeFilter(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	persistence := mocks.NewMockFeedRepository(ctrl)
	broadcast := mocks.NewMockBroadcastService(ctrl)
	users := mocks.NewMockUserRepository(ctrl)

	persisted := []domain.FeedEvent{{ID: "evt-1", EventType: domain.EventTypeUpload}}
	persistence.EXPECT().GetFeed(ctx, domain.EventTypeUpload, int64(10), int64(0)).Return(persisted, int64(1), nil)

	uc := NewFeedUseCase(persistence, users, broadcast, logger.New("info", false))
	entries, total, err := uc.GetFeed(ctx, domain.EventTypeUpload, 10, 0)

	require.NoError(t, err)
	require.Equal(t, persisted, entries)
	require.Equal(t, int64(1), total)
}
