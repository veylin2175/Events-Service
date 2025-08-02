package createEvent_test

import (
	"Events-Service/internal/http-server/handlers/event/createEvent"
	"bytes"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"Events-Service/internal/http-server/handlers/event/createEvent/mocks"
	"Events-Service/internal/storage"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNew_Success(t *testing.T) {
	mockService := new(mocks.CreateEvent)
	mockService.On("SaveEvent", mock.AnythingOfType("int64"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).
		Return(int64(42), nil).Once()

	requestBody := createEvent.Request{
		UserId: 1,
		Date:   "2025-08-05",
		Text:   "Test event",
	}
	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/events", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	testLogger := slog.New(slog.NewTextHandler(bytes.NewBuffer(nil), nil))
	handler := createEvent.New(testLogger, mockService)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var resp createEvent.Response
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, int64(42), resp.EventId)

	mockService.AssertExpectations(t)
}

func TestNew_EventExists(t *testing.T) {
	mockService := new(mocks.CreateEvent)
	mockService.On("SaveEvent", mock.Anything, mock.Anything, mock.Anything).
		Return(int64(0), storage.ErrEventExists).Once()

	requestBody := createEvent.Request{
		UserId: 1,
		Date:   "2025-08-05",
		Text:   "Test event",
	}
	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/events", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	testLogger := slog.New(slog.NewTextHandler(bytes.NewBuffer(nil), nil))
	handler := createEvent.New(testLogger, mockService)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusServiceUnavailable, rr.Code)

	mockService.AssertExpectations(t)
}

func TestNew_InternalServerError(t *testing.T) {
	mockService := new(mocks.CreateEvent)
	mockService.On("SaveEvent", mock.Anything, mock.Anything, mock.Anything).
		Return(int64(0), errors.New("database connection failed")).Once()

	requestBody := createEvent.Request{
		UserId: 1,
		Date:   "2025-08-05",
		Text:   "Test event",
	}
	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/events", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	testLogger := slog.New(slog.NewTextHandler(bytes.NewBuffer(nil), nil))
	handler := createEvent.New(testLogger, mockService)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	mockService.AssertExpectations(t)
}

func TestNew_ValidationError(t *testing.T) {
	mockService := new(mocks.CreateEvent)

	requestBody := createEvent.Request{
		UserId: 0,
		Date:   "",
		Text:   "",
	}
	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/events", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	testLogger := slog.New(slog.NewTextHandler(bytes.NewBuffer(nil), nil))
	handler := createEvent.New(testLogger, mockService)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	mockService.AssertNotCalled(t, "SaveEvent")
}
