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
	t.Parallel()

	tests := []struct {
		name       string
		req        PublishEventRequest
		setupMocks func(context.Context, *mocks.MockFeedRepository, *mocks.MockUserRepository, *mocks.MockBroadcastService) *domain.FeedEvent
		assertErr  func(*testing.T, error)
	}{
		{
			name: "when request is valid should publish event",
			req:  PublishEventRequest{EventType: domain.EventTypeNotification, Content: "deployment finished"},
			setupMocks: func(
				ctx context.Context,
				persistence *mocks.MockFeedRepository,
				users *mocks.MockUserRepository,
				broadcast *mocks.MockBroadcastService,
			) *domain.FeedEvent {
				created := &domain.FeedEvent{
					ID:        "evt-1",
					UserID:    "user-1",
					Username:  "alice",
					EventType: domain.EventTypeNotification,
					Content:   "deployment finished",
					CreatedAt: time.Now(),
				}

				users.EXPECT().GetByIDs(ctx, []string{"user-1"}).Return(map[string]string{"user-1": "alice"}, nil).Times(1)
				persistence.EXPECT().
					CreateEvent(ctx, domain.FeedEvent{
						UserID:    "user-1",
						Username:  "alice",
						EventType: domain.EventTypeNotification,
						Content:   "deployment finished",
					}).
					Return(created, nil).
					Times(1)
				broadcast.EXPECT().BroadcastEvent(ctx, created).Return(nil).Times(1)

				return created
			},
			assertErr: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		{
			name: "when user lookup fails should return wrapped error",
			req:  PublishEventRequest{EventType: domain.EventTypeNotification, Content: "deployment finished"},
			setupMocks: func(
				ctx context.Context,
				persistence *mocks.MockFeedRepository,
				users *mocks.MockUserRepository,
				broadcast *mocks.MockBroadcastService,
			) *domain.FeedEvent {
				repoErr := errors.New("user lookup failed")

				users.EXPECT().GetByIDs(ctx, []string{"user-1"}).Return(nil, repoErr).Times(1)
				persistence.EXPECT().CreateEvent(gomock.Any(), gomock.Any()).Times(0)
				broadcast.EXPECT().BroadcastEvent(gomock.Any(), gomock.Any()).Times(0)

				return nil
			},
			assertErr: func(t *testing.T, err error) {
				require.Error(t, err)
				require.ErrorContains(t, err, "failed to resolve event publisher")
			},
		},
		{
			name: "when persistence fails should return wrapped error",
			req:  PublishEventRequest{EventType: domain.EventTypeNotification, Content: "deployment finished"},
			setupMocks: func(
				ctx context.Context,
				persistence *mocks.MockFeedRepository,
				users *mocks.MockUserRepository,
				broadcast *mocks.MockBroadcastService,
			) *domain.FeedEvent {
				repoErr := errors.New("db error")

				users.EXPECT().GetByIDs(ctx, []string{"user-1"}).Return(map[string]string{"user-1": "alice"}, nil).Times(1)
				persistence.EXPECT().
					CreateEvent(ctx, domain.FeedEvent{
						UserID:    "user-1",
						Username:  "alice",
						EventType: domain.EventTypeNotification,
						Content:   "deployment finished",
					}).
					Return(nil, repoErr).
					Times(1)
				broadcast.EXPECT().BroadcastEvent(gomock.Any(), gomock.Any()).Times(0)

				return nil
			},
			assertErr: func(t *testing.T, err error) {
				require.Error(t, err)
				require.ErrorContains(t, err, "failed to publish event")
			},
		},
		{
			name: "when broadcast fails should still return created event",
			req:  PublishEventRequest{EventType: domain.EventTypeNotification, Content: "deployment finished"},
			setupMocks: func(
				ctx context.Context,
				persistence *mocks.MockFeedRepository,
				users *mocks.MockUserRepository,
				broadcast *mocks.MockBroadcastService,
			) *domain.FeedEvent {
				created := &domain.FeedEvent{
					ID:        "evt-2",
					UserID:    "user-1",
					Username:  "alice",
					EventType: domain.EventTypeNotification,
					Content:   "deployment finished",
					CreatedAt: time.Now(),
				}

				users.EXPECT().GetByIDs(ctx, []string{"user-1"}).Return(map[string]string{"user-1": "alice"}, nil).Times(1)
				persistence.EXPECT().
					CreateEvent(ctx, domain.FeedEvent{
						UserID:    "user-1",
						Username:  "alice",
						EventType: domain.EventTypeNotification,
						Content:   "deployment finished",
					}).
					Return(created, nil).
					Times(1)
				broadcast.EXPECT().BroadcastEvent(ctx, created).Return(errors.New("broadcast failed")).Times(1)

				return created
			},
			assertErr: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		{
			name: "when event type is invalid should return validation error",
			req:  PublishEventRequest{EventType: "freeform", Content: "deployment finished"},
			setupMocks: func(
				_ context.Context,
				persistence *mocks.MockFeedRepository,
				users *mocks.MockUserRepository,
				broadcast *mocks.MockBroadcastService,
			) *domain.FeedEvent {
				users.EXPECT().GetByIDs(gomock.Any(), gomock.Any()).Times(0)
				persistence.EXPECT().CreateEvent(gomock.Any(), gomock.Any()).Times(0)
				broadcast.EXPECT().BroadcastEvent(gomock.Any(), gomock.Any()).Times(0)

				return nil
			},
			assertErr: func(t *testing.T, err error) {
				require.Error(t, err)

				var invalidEventTypeErr *domain.InvalidEventTypeError
				require.ErrorAs(t, err, &invalidEventTypeErr)
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
			users := mocks.NewMockUserRepository(ctrl)
			broadcast := mocks.NewMockBroadcastService(ctrl)
			expected := tt.setupMocks(ctx, persistence, users, broadcast)

			uc := NewEventPublisherUseCase(persistence, users, broadcast, logger.New("info", false))

			// Act
			result, err := uc.PublishEvent(ctx, "user-1", tt.req)

			// Assert
			tt.assertErr(t, err)
			require.Equal(t, expected, result)
		})
	}
}
