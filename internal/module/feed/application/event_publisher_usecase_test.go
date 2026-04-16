package application

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"real-time-activity-feed/internal/module/feed/domain"
	"real-time-activity-feed/internal/module/feed/infrastructure/mocks"
	"real-time-activity-feed/internal/shared/logger"
)

func TestEventPublisherUseCase_PublishEvent(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	persistence := mocks.NewMockFeedRepository(ctrl)
	users := mocks.NewMockUserRepository(ctrl)
	broadcast := mocks.NewMockBroadcastService(ctrl)

	req := PublishEventRequest{EventType: "notification", Content: "deployment finished"}
	created := &domain.FeedEvent{
		ID: "evt-1", UserID: "user-1", Username: "alice", EventType: req.EventType, Content: req.Content, CreatedAt: time.Now(),
	}

	users.EXPECT().GetByIDs(ctx, []string{"user-1"}).Return(map[string]string{"user-1": "alice"}, nil)
	persistence.EXPECT().CreateEvent(ctx, gomock.Any()).Return(created, nil)
	broadcast.EXPECT().BroadcastEvent(ctx, created).Return(nil)

	uc := NewEventPublisherUseCase(persistence, users, broadcast, logger.New("info", false))
	result, err := uc.PublishEvent(ctx, "user-1", req)

	require.NoError(t, err)
	require.Equal(t, created, result)
}

func TestEventPublisherUseCase_PublishEvent_WhenPersistenceFails(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	persistence := mocks.NewMockFeedRepository(ctrl)
	users := mocks.NewMockUserRepository(ctrl)
	broadcast := mocks.NewMockBroadcastService(ctrl)

	users.EXPECT().GetByIDs(ctx, []string{"user-1"}).Return(map[string]string{"user-1": "alice"}, nil)
	persistence.EXPECT().CreateEvent(ctx, gomock.Any()).Return(nil, errors.New("db error"))

	uc := NewEventPublisherUseCase(persistence, users, broadcast, logger.New("info", false))
	_, err := uc.PublishEvent(ctx, "user-1", PublishEventRequest{EventType: "notification", Content: "deployment finished"})

	require.Error(t, err)
}
