// Package application provides use cases for the feed module.
package application

import (
	"context"
	"fmt"

	"real-time-activity-feed/internal/module/feed/domain"
	"real-time-activity-feed/internal/shared/logger"
)

//go:generate mockgen -destination=../adapters/mocks/event_publisher_usecase_mock.go -package=mocks real-time-activity-feed/internal/module/feed/application EventPublisherUseCase

// EventPublisherUseCase publishes new activity feed events.
type EventPublisherUseCase interface {
	PublishEvent(ctx context.Context, userID string, req PublishEventRequest) (*domain.FeedEvent, error)
}

type eventPublisherUseCase struct {
	persistenceRepo  FeedRepository
	userRepo         UserRepository
	broadcastService BroadcastService
	logger           *logger.Logger
}

// NewEventPublisherUseCase creates the event publish use case.
//
//nolint:revive // Returning the concrete type preserves the module's existing construction pattern.
func NewEventPublisherUseCase(
	persistenceRepo FeedRepository,
	userRepo UserRepository,
	broadcastService BroadcastService,
	l *logger.Logger,
) *eventPublisherUseCase {
	return &eventPublisherUseCase{
		persistenceRepo:  persistenceRepo,
		userRepo:         userRepo,
		broadcastService: broadcastService,
		logger:           l,
	}
}

// PublishEventRequest represents a request to publish a new feed event.
type PublishEventRequest struct {
	EventType string `json:"event_type" validate:"required,min=2,max=50" example:"notification"`
	Content   string `json:"content" validate:"required,min=3,max=255" example:"Alice uploaded a new report"`
}

// PublishEvent publishes a new feed event.
func (uc *eventPublisherUseCase) PublishEvent(ctx context.Context, userID string, req PublishEventRequest) (*domain.FeedEvent, error) {
	if err := domain.ValidateEventType(req.EventType); err != nil {
		return nil, err
	}

	usernames, err := uc.userRepo.GetByIDs(ctx, []string{userID})
	if err != nil {
		return nil, fmt.Errorf("failed to resolve event publisher: %w", err)
	}

	entry := domain.FeedEvent{
		UserID:    userID,
		Username:  usernames[userID],
		EventType: req.EventType,
		Content:   req.Content,
	}

	created, err := uc.persistenceRepo.CreateEvent(ctx, entry)
	if err != nil {
		uc.logger.Errorf(ctx, "Failed to persist feed event: %v", err)
		return nil, fmt.Errorf("failed to publish event: %w", err)
	}

	if err := uc.broadcastService.BroadcastEvent(ctx, created); err != nil {
		uc.logger.Warnf(ctx, "Failed to broadcast feed event: %v", err)
	}

	uc.logger.Infof(ctx, "Feed event published: user=%s type=%s", userID, req.EventType)
	return created, nil
}
