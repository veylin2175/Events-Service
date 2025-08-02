package deleteEvent_test

import (
	"Events-Service/internal/http-server/handlers/event/deleteEvent"
	"bytes"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"Events-Service/internal/http-server/handlers/event/deleteEvent/mocks"
	"Events-Service/internal/storage"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNew_Success(t *testing.T) {
	mockService := new(mocks.DeleteEvent)

	mockService.On("DeleteEvent", mock.AnythingOfType("int64"), mock.AnythingOfType("int64")).
		Return(nil).Once()

	requestBody := deleteEvent.Request{
		UserId:  1,
		EventId: 101,
	}
	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodDelete, "/events", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	testLogger := slog.New(slog.NewTextHandler(bytes.NewBuffer(nil), nil))
	handler := deleteEvent.New(testLogger, mockService)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	mockService.AssertExpectations(t)
}

func TestNew_EventNotFound(t *testing.T) {
	mockService := new(mocks.DeleteEvent)

	mockService.On("DeleteEvent", mock.AnythingOfType("int64"), mock.AnythingOfType("int64")).
		Return(storage.ErrEventNotFound).Once()

	requestBody := deleteEvent.Request{
		UserId:  1,
		EventId: 999,
	}
	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodDelete, "/events", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	testLogger := slog.New(slog.NewTextHandler(bytes.NewBuffer(nil), nil))
	handler := deleteEvent.New(testLogger, mockService)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusServiceUnavailable, rr.Code)

	mockService.AssertExpectations(t)
}

func TestNew_InternalServerError(t *testing.T) {
	mockService := new(mocks.DeleteEvent)

	mockService.On("DeleteEvent", mock.AnythingOfType("int64"), mock.AnythingOfType("int64")).
		Return(errors.New("database connection failed")).Once()

	requestBody := deleteEvent.Request{
		UserId:  1,
		EventId: 101,
	}
	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodDelete, "/events", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	testLogger := slog.New(slog.NewTextHandler(bytes.NewBuffer(nil), nil))
	handler := deleteEvent.New(testLogger, mockService)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	mockService.AssertExpectations(t)
}

func TestNew_ValidationError(t *testing.T) {
	mockService := new(mocks.DeleteEvent)

	requestBody := deleteEvent.Request{
		UserId:  0,
		EventId: 0,
	}
	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodDelete, "/events", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	testLogger := slog.New(slog.NewTextHandler(bytes.NewBuffer(nil), nil))
	handler := deleteEvent.New(testLogger, mockService)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	mockService.AssertNotCalled(t, "DeleteEvent")
}
