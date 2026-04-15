// Package v1 provides REST API v1 handlers for the feed module.
package v1

import (
	"encoding/json"
	"fmt"
	"time"

	"real-time-activity-feed/internal/module/feed/application"
	"real-time-activity-feed/internal/module/feed/domain"
	"real-time-activity-feed/internal/shared/logger"
	"real-time-activity-feed/internal/shared/middleware"
	"real-time-activity-feed/internal/shared/request"
	"real-time-activity-feed/internal/shared/response"
	"real-time-activity-feed/internal/shared/validator"

	"github.com/gin-gonic/gin"
)

const keepAliveInterval = 15 * time.Second

// FeedHandler exposes the activity feed API while preserving the existing module package.
type FeedHandler struct {
	feedUseCase  application.FeedUseCase
	eventUseCase application.EventPublisherUseCase
	logger       *logger.Logger
}

// NewFeedHandler creates an HTTP handler for the activity feed API.
func NewFeedHandler(
	feedUseCase application.FeedUseCase,
	eventUseCase application.EventPublisherUseCase,
	l *logger.Logger,
) *FeedHandler {
	return &FeedHandler{
		feedUseCase:  feedUseCase,
		eventUseCase: eventUseCase,
		logger:       l,
	}
}

// GetFeed handles `GET /feed`.
func (h *FeedHandler) GetFeed(c *gin.Context) {
	var pagination request.Pagination
	if err := c.ShouldBindQuery(&pagination); err != nil {
		apiErr := toAPIError(validator.Validate(pagination))
		response.Error(c, apiErr)
		return
	}

	if err := validator.Validate(pagination); err != nil {
		response.Error(c, toAPIError(err))
		return
	}

	normalized := pagination.Normalize()
	entries, total, err := h.feedUseCase.GetFeed(c.Request.Context(), normalized.GetLimit(), normalized.GetOffset())
	if err != nil {
		h.logger.Err(c.Request.Context(), err).Msg("Feed request error")
		response.Error(c, toAPIError(err))
		return
	}

	meta := response.NewPagination(normalized.GetOffset(), normalized.GetLimit(), total)
	response.SuccessWithMeta(c, entries, "Activity feed retrieved successfully", meta)
}

// StreamFeed handles `GET /feed/stream` using SSE.
func (h *FeedHandler) StreamFeed(c *gin.Context) {
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")

	ctx := c.Request.Context()
	updateCh, err := h.feedUseCase.SubscribeToFeedEvents(ctx)
	if err != nil {
		closedCh := make(chan *domain.FeedEvent)
		close(closedCh)
		updateCh = closedCh
	}

	ticker := time.NewTicker(keepAliveInterval)
	defer ticker.Stop()

	notify := c.Writer.CloseNotify()
	for {
		select {
		case <-notify:
			return
		case entry, ok := <-updateCh:
			if !ok {
				return
			}

			resp := response.Response{
				Success: true,
				Data:    entry,
				Message: "Feed event received",
			}
			messageBytes, _ := json.Marshal(resp)
			_, _ = fmt.Fprintf(c.Writer, "data: %s\n\n", messageBytes)
			c.Writer.Flush()
		case <-ticker.C:
			_, _ = fmt.Fprintf(c.Writer, ": keep-alive\n\n")
			c.Writer.Flush()
		}
	}
}

// PublishEvent handles `POST /events`.
func (h *FeedHandler) PublishEvent(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Error(c, response.NewInternalError("An unexpected error occurred"))
		return
	}

	var req application.PublishEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, toAPIError(validator.Validate(req)))
		return
	}

	if err := validator.Validate(req); err != nil {
		response.Error(c, toAPIError(err))
		return
	}

	entry, err := h.eventUseCase.PublishEvent(c.Request.Context(), userID, req)
	if err != nil {
		h.logger.Err(c.Request.Context(), err).Msg("Publish event error")
		response.Error(c, toAPIError(err))
		return
	}

	response.Success(c, entry, "Event published successfully")
}

// RegisterPublicRoutes registers public feed endpoints.
func (h *FeedHandler) RegisterPublicRoutes(router *gin.RouterGroup) {
	feed := router.Group("/feed")
	{
		feed.GET("", h.GetFeed)
		feed.GET("/stream", h.StreamFeed)
	}
}

// RegisterProtectedRoutes registers authenticated event publishing endpoints.
func (h *FeedHandler) RegisterProtectedRoutes(router *gin.RouterGroup) {
	events := router.Group("/events")
	{
		events.POST("", h.PublishEvent)
	}
}
