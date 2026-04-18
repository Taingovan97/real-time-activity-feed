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
	GetFeed(ctx context.Context, eventType, query string, limit, offset int64) ([]domain.FeedEvent, int64, error)
	SubscribeToFeedEvents(ctx context.Context) (<-chan *domain.FeedEvent, error)
}

type feedUseCase struct {
	persistenceRepo  FeedRepository
	cache            FeedCache
	broadcastService BroadcastService
	logger           *logger.Logger
}

// NewFeedUseCase creates the feed read use case.
//
//nolint:revive // Returning the concrete type preserves the module's existing construction pattern.
func NewFeedUseCase(
	persistenceRepo FeedRepository,
	cache FeedCache,
	_ UserRepository,
	broadcastService BroadcastService,
	l *logger.Logger,
) *feedUseCase {
	return &feedUseCase{
		persistenceRepo:  persistenceRepo,
		cache:            cache,
		broadcastService: broadcastService,
		logger:           l,
	}
}

// GetFeed returns feed items.
func (uc *feedUseCase) GetFeed(ctx context.Context, eventType, query string, limit, offset int64) ([]domain.FeedEvent, int64, error) {
	cacheable := query == "" && offset == 0 && limit == 10

	if cacheable {
		entries, total, found, err := uc.cache.GetFeed(ctx, eventType, limit, offset)
		if err != nil {
			uc.logger.Warnf(ctx, "Failed to read feed cache: %v", err)
		} else if found {
			return entries, total, nil
		}
	}

	entries, total, err := uc.persistenceRepo.GetFeed(ctx, eventType, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to retrieve feed: %w", err)
	}

	if cacheable {
		if err := uc.cache.SetFeed(ctx, eventType, limit, offset, entries, total); err != nil {
			uc.logger.Warnf(ctx, "Failed to write feed cache: %v", err)
		}
	}

	return entries, total, nil
}

// SubscribeToFeedEvents subscribes callers to live feed events.
func (uc *feedUseCase) SubscribeToFeedEvents(ctx context.Context) (<-chan *domain.FeedEvent, error) {
	return uc.broadcastService.SubscribeToFeedEvents(ctx)
}
