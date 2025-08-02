package updateEvent_test

import (
	"Events-Service/internal/http-server/handlers/event/updateEvent"
	"bytes"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"Events-Service/internal/http-server/handlers/event/updateEvent/mocks"
	"Events-Service/internal/storage"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNew_Success(t *testing.T) {
	mockService := new(mocks.UpdateEvent)

	mockService.On("UpdateEvent", mock.AnythingOfType("int64"), mock.AnythingOfType("int64"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).
		Return(nil).Once()

	requestBody := updateEvent.Request{
		UserId:  1,
		EventId: 101,
		Date:    "2025-08-05",
		Text:    "Updated event",
	}
	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPut, "/events", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	testLogger := slog.New(slog.NewTextHandler(bytes.NewBuffer(nil), nil))
	handler := updateEvent.New(testLogger, mockService)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	mockService.AssertExpectations(t)
}

func TestNew_EventNotFound(t *testing.T) {
	mockService := new(mocks.UpdateEvent)

	mockService.On("UpdateEvent", mock.AnythingOfType("int64"), mock.AnythingOfType("int64"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).
		Return(storage.ErrEventNotFound).Once()

	requestBody := updateEvent.Request{
		UserId:  1,
		EventId: 999,
		Date:    "2025-08-05",
		Text:    "Updated event",
	}
	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPut, "/events", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	testLogger := slog.New(slog.NewTextHandler(bytes.NewBuffer(nil), nil))
	handler := updateEvent.New(testLogger, mockService)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusServiceUnavailable, rr.Code)

	mockService.AssertExpectations(t)
}

func TestNew_InternalServerError(t *testing.T) {
	mockService := new(mocks.UpdateEvent)

	mockService.On("UpdateEvent", mock.AnythingOfType("int64"), mock.AnythingOfType("int64"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).
		Return(errors.New("database connection failed")).Once()

	requestBody := updateEvent.Request{
		UserId:  1,
		EventId: 101,
		Date:    "2025-08-05",
		Text:    "Updated event",
	}
	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPut, "/events", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	testLogger := slog.New(slog.NewTextHandler(bytes.NewBuffer(nil), nil))
	handler := updateEvent.New(testLogger, mockService)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	mockService.AssertExpectations(t)
}

func TestNew_ValidationError(t *testing.T) {
	mockService := new(mocks.UpdateEvent)

	requestBody := updateEvent.Request{
		UserId:  0,
		EventId: 0,
		Date:    "",
		Text:    "",
	}
	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPut, "/events", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	testLogger := slog.New(slog.NewTextHandler(bytes.NewBuffer(nil), nil))
	handler := updateEvent.New(testLogger, mockService)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	mockService.AssertNotCalled(t, "UpdateEvent")
}
