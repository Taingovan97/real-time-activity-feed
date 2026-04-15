// Package application provides use cases for the feed module.
package application

//go:generate mockgen -destination=../infrastructure/mocks/broadcast_service_mock.go -package=mocks real-time-activity-feed/internal/module/feed/application BroadcastService

import (
	"context"

	"real-time-activity-feed/internal/module/feed/domain"
)

// BroadcastService defines the interface for broadcasting feed updates.
type BroadcastService interface {
	BroadcastEvent(ctx context.Context, entry *domain.FeedEvent) error
	SubscribeToFeedEvents(ctx context.Context) (<-chan *domain.FeedEvent, error)
}
