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

func TestFeedUseCase_GetFeed(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		eventType    string
		query        string
		limit        int64
		offset       int64
		setupMocks   func(context.Context, *mocks.MockFeedRepository)
		assertResult func(*testing.T, []domain.FeedEvent, int64, error)
	}{
		{
			name:   "when repository returns entries should return feed",
			limit:  10,
			offset: 0,
			setupMocks: func(ctx context.Context, persistence *mocks.MockFeedRepository) {
				persisted := []domain.FeedEvent{{ID: "evt-1"}}
				persistence.EXPECT().GetFeed(ctx, "", "", int64(10), int64(0)).Return(persisted, int64(1), nil).Times(1)
			},
			assertResult: func(t *testing.T, entries []domain.FeedEvent, total int64, err error) {
				require.NoError(t, err)
				require.Equal(t, []domain.FeedEvent{{ID: "evt-1"}}, entries)
				require.Equal(t, int64(1), total)
			},
		},
		{
			name:   "when repository fails should return wrapped error",
			limit:  10,
			offset: 0,
			setupMocks: func(ctx context.Context, persistence *mocks.MockFeedRepository) {
				persistence.EXPECT().GetFeed(ctx, "", "", int64(10), int64(0)).Return(nil, int64(0), errors.New("db down")).Times(1)
			},
			assertResult: func(t *testing.T, entries []domain.FeedEvent, total int64, err error) {
				require.Error(t, err)
				require.ErrorContains(t, err, "failed to retrieve feed")
				require.Nil(t, entries)
				require.Zero(t, total)
			},
		},
		{
			name:      "when event type filter is provided should pass it through",
			eventType: domain.EventTypeUpload,
			limit:     10,
			offset:    0,
			setupMocks: func(ctx context.Context, persistence *mocks.MockFeedRepository) {
				persisted := []domain.FeedEvent{{ID: "evt-1", EventType: domain.EventTypeUpload}}
				persistence.EXPECT().GetFeed(ctx, domain.EventTypeUpload, "", int64(10), int64(0)).Return(persisted, int64(1), nil).Times(1)
			},
			assertResult: func(t *testing.T, entries []domain.FeedEvent, total int64, err error) {
				require.NoError(t, err)
				require.Equal(t, []domain.FeedEvent{{ID: "evt-1", EventType: domain.EventTypeUpload}}, entries)
				require.Equal(t, int64(1), total)
			},
		},
		{
			name:      "when event type and query are provided should pass both through",
			eventType: domain.EventTypeUpload,
			query:     "report",
			limit:     10,
			offset:    0,
			setupMocks: func(ctx context.Context, persistence *mocks.MockFeedRepository) {
				persisted := []domain.FeedEvent{{ID: "evt-1", EventType: domain.EventTypeUpload, Content: "uploaded report"}}
				persistence.EXPECT().GetFeed(ctx, domain.EventTypeUpload, "report", int64(10), int64(0)).Return(persisted, int64(1), nil).Times(1)
			},
			assertResult: func(t *testing.T, entries []domain.FeedEvent, total int64, err error) {
				require.NoError(t, err)
				require.Equal(t, []domain.FeedEvent{{ID: "evt-1", EventType: domain.EventTypeUpload, Content: "uploaded report"}}, entries)
				require.Equal(t, int64(1), total)
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			ctx := context.Background()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			persistence := mocks.NewMockFeedRepository(ctrl)
			broadcast := mocks.NewMockBroadcastService(ctrl)
			users := mocks.NewMockUserRepository(ctrl)
			tt.setupMocks(ctx, persistence)

			uc := NewFeedUseCase(persistence, users, broadcast, logger.New("info", false))

			// Act
			entries, total, err := uc.GetFeed(ctx, tt.eventType, tt.query, tt.limit, tt.offset)

			// Assert
			tt.assertResult(t, entries, total, err)
		})
	}
}

func TestFeedUseCase_SubscribeToFeedEvents(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		setupMocks   func(context.Context, *mocks.MockBroadcastService) <-chan *domain.FeedEvent
		assertResult func(*testing.T, <-chan *domain.FeedEvent, error)
	}{
		{
			name: "when broadcast service subscribes should return channel",
			setupMocks: func(ctx context.Context, broadcast *mocks.MockBroadcastService) <-chan *domain.FeedEvent {
				events := make(chan *domain.FeedEvent, 1)
				broadcast.EXPECT().SubscribeToFeedEvents(ctx).Return(events, nil).Times(1)
				return events
			},
			assertResult: func(t *testing.T, events <-chan *domain.FeedEvent, err error) {
				require.NoError(t, err)
				require.NotNil(t, events)
			},
		},
		{
			name: "when broadcast service fails should return error",
			setupMocks: func(ctx context.Context, broadcast *mocks.MockBroadcastService) <-chan *domain.FeedEvent {
				broadcast.EXPECT().SubscribeToFeedEvents(ctx).Return(nil, errors.New("subscribe failed")).Times(1)
				return nil
			},
			assertResult: func(t *testing.T, events <-chan *domain.FeedEvent, err error) {
				require.Error(t, err)
				require.ErrorContains(t, err, "subscribe failed")
				require.Nil(t, events)
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			ctx := context.Background()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			persistence := mocks.NewMockFeedRepository(ctrl)
			broadcast := mocks.NewMockBroadcastService(ctrl)
			users := mocks.NewMockUserRepository(ctrl)
			expectedEvents := tt.setupMocks(ctx, broadcast)

			uc := NewFeedUseCase(persistence, users, broadcast, logger.New("info", false))

			// Act
			events, err := uc.SubscribeToFeedEvents(ctx)

			// Assert
			tt.assertResult(t, events, err)
			require.Equal(t, expectedEvents, events)
		})
	}
}
