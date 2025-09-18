package main_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"net"
	"net/http"
	"os"
	"sync"
	"testing"
	"time"

	"Events-Service/internal/config"
	"Events-Service/internal/http-server/handlers/event/createEvent"
	"Events-Service/internal/http-server/handlers/user"
	"Events-Service/internal/http-server/middleware/mwlogger"
	"Events-Service/internal/lib/logger/sl"
	"Events-Service/internal/storage/postgres"
)

func setupTestServer() (string, *sync.WaitGroup, func(), error) {
	cfg := &config.Config{
		HTTPServer: config.HTTPServer{
			Address:     ":0",
			Timeout:     10 * time.Second,
			IdleTimeout: 60 * time.Second,
		},
		Database: config.Database{
			Host:     "localhost",
			Port:     5432,
			User:     "postgres",
			Password: "3356",
			DBName:   "events_service",
			SSLMode:  "disable",
		},
	}

	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	db, err := postgres.InitDB(cfg)
	if err != nil {
		return "", nil, nil, fmt.Errorf("failed to init test storage: %w", err)
	}

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(mwlogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Post("/create_user", user.New(log, db))
	router.Post("/create_event", createEvent.New(log, db))

	srv := &http.Server{
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	listener, err := net.Listen("tcp", cfg.HTTPServer.Address)
	if err != nil {
		db.Close()
		return "", nil, nil, fmt.Errorf("failed to create listener: %w", err)
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		if err = srv.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("failed to start test server", sl.Err(err))
		}
	}()

	testServerAddr = listener.Addr().String()

	teardown := func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err = srv.Shutdown(ctx); err != nil {
			log.Error("failed to stop test server", sl.Err(err))
		}
		db.Close()
	}

	return testServerAddr, wg, teardown, nil
}

var testServerAddr string
var teardownServer func()

func TestMain(m *testing.M) {
	var wg *sync.WaitGroup
	var err error

	testServerAddr, wg, teardownServer, err = setupTestServer()
	if err != nil {
		slog.Error("failed to setup test server", sl.Err(err))
		os.Exit(1)
	}

	code := m.Run()

	teardownServer()
	wg.Wait()

	os.Exit(code)
}

func TestCreateUserAndEvent(t *testing.T) {
	userReqBody := `{}`
	userReq, err := http.NewRequest("POST", "http://"+testServerAddr+"/create_user", bytes.NewBufferString(userReqBody))
	assert.NoError(t, err)
	userReq.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(userReq)
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var userResp user.Response
	err = json.NewDecoder(resp.Body).Decode(&userResp)
	assert.NoError(t, err)
	assert.True(t, userResp.UserId > 0)

	eventReqBody := createEvent.Request{
		UserId: userResp.UserId,
		Date:   "2025-12-25",
		Text:   "Christmas party",
	}
	eventBody, _ := json.Marshal(eventReqBody)
	eventReq, err := http.NewRequest("POST", "http://"+testServerAddr+"/create_event", bytes.NewBuffer(eventBody))
	assert.NoError(t, err)
	eventReq.Header.Set("Content-Type", "application/json")

	resp, err = http.DefaultClient.Do(eventReq)
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var eventResp createEvent.Response
	err = json.NewDecoder(resp.Body).Decode(&eventResp)
	assert.NoError(t, err)
	assert.True(t, eventResp.EventId > 0)
}

func TestCreateUser_InvalidJson(t *testing.T) {
	req, err := http.NewRequest("POST", "http://"+testServerAddr+"/create_user", bytes.NewBufferString("invalid json"))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}
