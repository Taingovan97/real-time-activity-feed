package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"real-time-activity-feed/internal/module/feed/adapters/mocks"
	"real-time-activity-feed/internal/module/feed/application"
	"real-time-activity-feed/internal/module/feed/domain"
	"real-time-activity-feed/internal/shared/logger"
	"real-time-activity-feed/internal/shared/response"
)

func TestFeedHandler_GetFeed(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFeed := mocks.NewMockFeedUseCase(ctrl)
	mockEvents := mocks.NewMockEventPublisherUseCase(ctrl)
	mockFeed.EXPECT().GetFeed(gomock.Any(), "", int64(10), int64(0)).Return([]domain.FeedEvent{{
		ID: "evt-1", UserID: "user-1", Username: "alice", EventType: "login", Content: "signed in", CreatedAt: time.Now(),
	}}, int64(1), nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/feed?limit=10&offset=0", nil)

	NewFeedHandler(mockFeed, mockEvents, logger.New("info", false)).GetFeed(c)

	require.Equal(t, http.StatusOK, w.Code)
	var body response.Response
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	require.True(t, body.Success)
	require.Equal(t, "Activity feed retrieved successfully", body.Message)
}

func TestFeedHandler_GetFeed_WithEventTypeFilter(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFeed := mocks.NewMockFeedUseCase(ctrl)
	mockEvents := mocks.NewMockEventPublisherUseCase(ctrl)
	mockFeed.EXPECT().GetFeed(gomock.Any(), domain.EventTypeUpload, int64(10), int64(0)).Return([]domain.FeedEvent{{
		ID: "evt-1", UserID: "user-1", Username: "alice", EventType: domain.EventTypeUpload, Content: "uploaded file", CreatedAt: time.Now(),
	}}, int64(1), nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/feed?limit=10&offset=0&event_type=upload", nil)

	NewFeedHandler(mockFeed, mockEvents, logger.New("info", false)).GetFeed(c)

	require.Equal(t, http.StatusOK, w.Code)
}

func TestFeedHandler_GetFeed_WhenEventTypeIsInvalid(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFeed := mocks.NewMockFeedUseCase(ctrl)
	mockEvents := mocks.NewMockEventPublisherUseCase(ctrl)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/feed?event_type=freeform", nil)

	NewFeedHandler(mockFeed, mockEvents, logger.New("info", false)).GetFeed(c)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestFeedHandler_ListEventTypes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFeed := mocks.NewMockFeedUseCase(ctrl)
	mockEvents := mocks.NewMockEventPublisherUseCase(ctrl)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/feed/event-types", nil)

	NewFeedHandler(mockFeed, mockEvents, logger.New("info", false)).ListEventTypes(c)

	require.Equal(t, http.StatusOK, w.Code)
	var body response.Response
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	require.True(t, body.Success)
	require.Equal(t, "Event types retrieved successfully", body.Message)
}

func TestFeedHandler_StreamFeed_WebSocket(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFeed := mocks.NewMockFeedUseCase(ctrl)
	mockEvents := mocks.NewMockEventPublisherUseCase(ctrl)
	updates := make(chan *domain.FeedEvent, 1)
	mockFeed.EXPECT().SubscribeToFeedEvents(gomock.Any()).DoAndReturn(
		func(_ context.Context) (<-chan *domain.FeedEvent, error) {
			return updates, nil
		},
	)

	router := gin.New()
	NewFeedHandler(mockFeed, mockEvents, logger.New("info", false)).RegisterPublicRoutes(&router.RouterGroup)
	server := httptest.NewServer(router)
	defer server.Close()

	wsURL := "ws" + server.URL[len("http"):] + "/feed/ws"
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, conn.Close())
	}()

	updates <- &domain.FeedEvent{
		ID: "evt-1", UserID: "user-1", Username: "alice", EventType: "login", Content: "signed in", CreatedAt: time.Now(),
	}

	_, message, err := conn.ReadMessage()
	require.NoError(t, err)

	var body response.Response
	require.NoError(t, json.Unmarshal(message, &body))
	require.True(t, body.Success)
	require.Equal(t, "Feed event received", body.Message)
}

func TestFeedHandler_PublishEvent(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFeed := mocks.NewMockFeedUseCase(ctrl)
	mockEvents := mocks.NewMockEventPublisherUseCase(ctrl)
	req := application.PublishEventRequest{EventType: "notification", Content: "hello world"}
	mockEvents.EXPECT().PublishEvent(gomock.Any(), "user-123", req).Return(&domain.FeedEvent{
		ID: "evt-1", UserID: "user-123", EventType: "notification", Content: "hello world",
	}, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/events", bytes.NewBufferString(`{"event_type":"notification","content":"hello world"}`))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("user_id", "user-123")

	NewFeedHandler(mockFeed, mockEvents, logger.New("info", false)).PublishEvent(c)

	require.Equal(t, http.StatusOK, w.Code)
}

func TestFeedHandler_PublishEvent_WhenUseCaseFails(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFeed := mocks.NewMockFeedUseCase(ctrl)
	mockEvents := mocks.NewMockEventPublisherUseCase(ctrl)
	mockEvents.EXPECT().PublishEvent(gomock.Any(), "user-123", gomock.Any()).Return(nil, errors.New("boom"))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/events", bytes.NewBufferString(`{"event_type":"notification","content":"hello world"}`))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("user_id", "user-123")

	NewFeedHandler(mockFeed, mockEvents, logger.New("info", false)).PublishEvent(c)

	require.Equal(t, http.StatusInternalServerError, w.Code)
}
