package application

//go:generate mockgen -destination=../infrastructure/mocks/repository_mock.go -package=mocks real-time-activity-feed/internal/module/feed/application UserRepository,FeedRepository,FeedCache

import (
	"context"

	"real-time-activity-feed/internal/module/feed/domain"
)

// UserRepository defines the interface for user data operations in the feed module.
type UserRepository interface {
	GetByIDs(ctx context.Context, userIDs []string) (map[string]string, error)
}

// FeedRepository stores feed events in PostgreSQL.
type FeedRepository interface {
	CreateEvent(ctx context.Context, entry domain.FeedEvent) (*domain.FeedEvent, error)
	GetFeed(ctx context.Context, eventType, query string, limit, offset int64) ([]domain.FeedEvent, int64, error)
}

// FeedCache stores hot feed pages in a non-authoritative cache.
type FeedCache interface {
	GetFeed(ctx context.Context, eventType string, limit, offset int64) ([]domain.FeedEvent, int64, bool, error)
	SetFeed(ctx context.Context, eventType string, limit, offset int64, entries []domain.FeedEvent, total int64) error
	InvalidateAfterPublish(ctx context.Context, eventType string) error
}
