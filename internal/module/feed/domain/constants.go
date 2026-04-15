// Package domain provides domain entities and constants for the activity feed module.
package domain

const (
	// RedisViewerUpdateTopic is the Redis pub/sub topic published with feed item updates.
	RedisViewerUpdateTopic = "feed:viewer:updates"

	// RedisFeedKey stores the recent activity feed items in Redis.
	RedisFeedKey = "feed:recent"

	// MaxFeedItems is the max number of feed items retained in Redis.
	MaxFeedItems = 1000
)
