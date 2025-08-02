package getEvents_test

import (
	"Events-Service/internal/http-server/handlers/event/getEvents"
	"Events-Service/internal/http-server/handlers/event/getEvents/mocks"
	"Events-Service/internal/models"
	"bytes"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestByWeek_Success(t *testing.T) {
	mockService := new(mocks.GetEvents)

	mockService.On("GetEventsByWeek", mock.AnythingOfType("int64"), mock.AnythingOfType("time.Time")).
		Return([]models.Event{
			{Date: "2025-08-02", Text: "Event A"},
			{Date: "2025-08-05", Text: "Event B"},
		}, nil).Once()

	requestBody := getEvents.Request{
		UserId: 1,
		Date:   "2025-08-05",
	}
	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/events/by-week", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	testLogger := slog.New(slog.NewTextHandler(bytes.NewBuffer(nil), nil))
	handler := getEvents.ByWeek(testLogger, mockService)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var resp getEvents.Response
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	assert.NoError(t, err)

	expectedEvents := []getEvents.EventResponse{
		{Date: "2025-08-02", Text: "Event A"},
		{Date: "2025-08-05", Text: "Event B"},
	}
	assert.Equal(t, expectedEvents, resp.Events)

	mockService.AssertExpectations(t)
}

func TestByWeek_ServiceError(t *testing.T) {
	mockService := new(mocks.GetEvents)

	mockService.On("GetEventsByWeek", mock.AnythingOfType("int64"), mock.AnythingOfType("time.Time")).
		Return(nil, errors.New("database error")).Once()

	requestBody := getEvents.Request{
		UserId: 1,
		Date:   "2025-08-05",
	}
	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/events/by-week", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	testLogger := slog.New(slog.NewTextHandler(bytes.NewBuffer(nil), nil))
	handler := getEvents.ByWeek(testLogger, mockService)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	mockService.AssertExpectations(t)
}

func TestByWeek_InvalidDate(t *testing.T) {
	mockService := new(mocks.GetEvents)

	requestBody := getEvents.Request{
		UserId: 1,
		Date:   "invalid-date",
	}
	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/events_for_week", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	testLogger := slog.New(slog.NewTextHandler(bytes.NewBuffer(nil), nil))
	handler := getEvents.ByWeek(testLogger, mockService)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	mockService.AssertNotCalled(t, "GetEventsByWeek")
}
