package v1

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
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
	mockFeed.EXPECT().GetFeed(gomock.Any(), int64(10), int64(0)).Return([]domain.FeedEvent{{
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
