// Package v1 provides REST API v1 handlers for the feed module.
package v1

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"real-time-activity-feed/internal/module/feed/application"
	"real-time-activity-feed/internal/shared/logger"
	"real-time-activity-feed/internal/shared/middleware"
	"real-time-activity-feed/internal/shared/request"
	"real-time-activity-feed/internal/shared/response"
	"real-time-activity-feed/internal/shared/validator"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const (
	websocketPingInterval = 30 * time.Second
	websocketWriteWait    = 10 * time.Second
	websocketPongWait     = 60 * time.Second
	maxMessageSize        = 1024
)

var feedUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(_ *http.Request) bool {
		return true
	},
}

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

// StreamFeed handles `GET /feed/ws` using WebSocket.
func (h *FeedHandler) StreamFeed(c *gin.Context) {
	ctx, cancel := context.WithCancel(c.Request.Context())
	defer cancel()

	updateCh, err := h.feedUseCase.SubscribeToFeedEvents(ctx)
	if err != nil {
		response.Error(c, response.NewInternalError("Unable to connect live feed"))
		return
	}

	conn, err := feedUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		h.logger.Warnf(ctx, "Failed to upgrade feed WebSocket connection: %v", err)
		return
	}
	defer func() {
		_ = conn.Close()
	}()

	conn.SetReadLimit(maxMessageSize)
	_ = conn.SetReadDeadline(time.Now().Add(websocketPongWait))
	conn.SetPongHandler(func(string) error {
		return conn.SetReadDeadline(time.Now().Add(websocketPongWait))
	})

	go func() {
		defer cancel()
		for {
			if _, _, readErr := conn.ReadMessage(); readErr != nil {
				return
			}
		}
	}()

	ticker := time.NewTicker(websocketPingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
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
			messageBytes, marshalErr := json.Marshal(resp)
			if marshalErr != nil {
				h.logger.Warnf(ctx, "Failed to marshal WebSocket feed event: %v", marshalErr)
				continue
			}
			_ = conn.SetWriteDeadline(time.Now().Add(websocketWriteWait))
			if writeErr := conn.WriteMessage(websocket.TextMessage, messageBytes); writeErr != nil {
				h.logger.Warnf(ctx, "Failed to write WebSocket feed event: %v", writeErr)
				return
			}
		case <-ticker.C:
			if pingErr := conn.WriteControl(
				websocket.PingMessage,
				[]byte("ping"),
				time.Now().Add(websocketWriteWait),
			); pingErr != nil {
				h.logger.Warnf(ctx, "Failed to ping WebSocket client: %v", pingErr)
				return
			}
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
		feed.GET("/ws", h.StreamFeed)
	}
}

// RegisterProtectedRoutes registers authenticated event publishing endpoints.
func (h *FeedHandler) RegisterProtectedRoutes(router *gin.RouterGroup) {
	events := router.Group("/events")
	{
		events.POST("", h.PublishEvent)
	}
}
