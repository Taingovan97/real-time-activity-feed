// Package broadcast provides broadcast service implementations for the feed module.
package broadcast

import (
	"context"
	"encoding/json"
	"fmt"

	"real-time-activity-feed/internal/module/feed/application"
	"real-time-activity-feed/internal/module/feed/domain"
	"real-time-activity-feed/internal/shared/logger"

	"github.com/redis/go-redis/v9"
)

// RedisBroadcastService implements Redis-backed fan-out for feed events.
type RedisBroadcastService struct {
	client      *redis.Client
	logger      *logger.Logger
	viewerTopic string
}

// NewRedisBroadcastService creates a Redis-backed feed broadcaster.
func NewRedisBroadcastService(
	client *redis.Client,
	logger *logger.Logger,
) application.BroadcastService {
	return &RedisBroadcastService{
		client:      client,
		logger:      logger,
		viewerTopic: domain.RedisViewerUpdateTopic,
	}
}

// BroadcastEvent publishes a single feed item to Redis pub/sub.
func (s *RedisBroadcastService) BroadcastEvent(ctx context.Context, entry *domain.FeedEvent) error {
	jsonData, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	return s.client.Publish(ctx, s.viewerTopic, jsonData).Err()
}

// SubscribeToFeedEvents subscribes to live feed items from Redis pub/sub.
func (s *RedisBroadcastService) SubscribeToFeedEvents(ctx context.Context) (<-chan *domain.FeedEvent, error) {
	pubsub := s.client.Subscribe(ctx, s.viewerTopic)
	ch := make(chan *domain.FeedEvent, 1)

	go func() {
		defer close(ch)
		defer func() {
			_ = pubsub.Close()
		}()

		for {
			select {
			case <-ctx.Done():
				return
			case msg := <-pubsub.Channel():
				if msg == nil {
					return
				}

				var entry domain.FeedEvent
				if err := json.Unmarshal([]byte(msg.Payload), &entry); err != nil {
					s.logger.Warnf(ctx, "Failed to unmarshal feed item: %v", err)
					continue
				}

				select {
				case ch <- &entry:
				case <-ctx.Done():
					return
				}
			}
		}
	}()

	return ch, nil
}
