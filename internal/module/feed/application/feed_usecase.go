// Package application provides use cases for the feed module.
package application

import (
	"context"
	"fmt"

	"real-time-activity-feed/internal/module/feed/domain"
	"real-time-activity-feed/internal/shared/logger"
)

//go:generate mockgen -destination=../adapters/mocks/feed_usecase_mock.go -package=mocks real-time-activity-feed/internal/module/feed/application FeedUseCase

// FeedUseCase serves the activity feed.
type FeedUseCase interface {
	GetFeed(ctx context.Context, limit, offset int64) ([]domain.FeedEvent, int64, error)
	SubscribeToFeedEvents(ctx context.Context) (<-chan *domain.FeedEvent, error)
}

type feedUseCase struct {
	cacheRepo        FeedCacheRepository
	persistenceRepo  FeedRepository
	broadcastService BroadcastService
	logger           *logger.Logger
}

// NewFeedUseCase creates the feed read use case.
//
//nolint:revive // Returning the concrete type preserves the module's existing construction pattern.
func NewFeedUseCase(
	cacheRepo FeedCacheRepository,
	persistenceRepo FeedRepository,
	_ UserRepository,
	broadcastService BroadcastService,
	l *logger.Logger,
) *feedUseCase {
	return &feedUseCase{
		cacheRepo:        cacheRepo,
		persistenceRepo:  persistenceRepo,
		broadcastService: broadcastService,
		logger:           l,
	}
}

// GetFeed returns feed items.
func (uc *feedUseCase) GetFeed(ctx context.Context, limit, offset int64) ([]domain.FeedEvent, int64, error) {
	entries, total, err := uc.cacheRepo.GetFeed(ctx, limit, offset)
	if err == nil && total > 0 {
		return entries, total, nil
	}

	if err != nil {
		uc.logger.Warnf(ctx, "Feed cache unavailable, falling back to PostgreSQL: %v", err)
	}

	entries, total, err = uc.persistenceRepo.GetFeed(ctx, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to retrieve feed: %w", err)
	}

	return entries, total, nil
}

// SubscribeToFeedEvents subscribes callers to live feed events.
func (uc *feedUseCase) SubscribeToFeedEvents(ctx context.Context) (<-chan *domain.FeedEvent, error) {
	return uc.broadcastService.SubscribeToFeedEvents(ctx)
}
