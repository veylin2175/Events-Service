package user_test

import (
	"Events-Service/internal/http-server/handlers/user"
	"bytes"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"Events-Service/internal/http-server/handlers/user/mocks"

	"github.com/stretchr/testify/assert"
)

func TestNew_Success(t *testing.T) {
	mockService := new(mocks.UserCreator)

	mockService.On("CreateUser").Return(int64(42), nil).Once()

	reqBody, _ := json.Marshal(user.Request{})
	req := httptest.NewRequest(http.MethodPost, "/user", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	testLogger := slog.New(slog.NewTextHandler(bytes.NewBuffer(nil), nil))
	handler := user.New(testLogger, mockService)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var resp user.Response
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, int64(42), resp.UserId)

	mockService.AssertExpectations(t)
}

func TestNew_InternalServerError(t *testing.T) {
	mockService := new(mocks.UserCreator)

	mockService.On("CreateUser").Return(int64(0), errors.New("database error")).Once()

	reqBody, _ := json.Marshal(user.Request{})
	req := httptest.NewRequest(http.MethodPost, "/user", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	testLogger := slog.New(slog.NewTextHandler(bytes.NewBuffer(nil), nil))
	handler := user.New(testLogger, mockService)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	mockService.AssertExpectations(t)
}

func TestNew_BadRequest(t *testing.T) {
	mockService := new(mocks.UserCreator)

	req := httptest.NewRequest(http.MethodPost, "/user", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	testLogger := slog.New(slog.NewTextHandler(bytes.NewBuffer(nil), nil))
	handler := user.New(testLogger, mockService)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	mockService.AssertNotCalled(t, "CreateUser")
}
