package main

import (
	"Events-Service/internal/config"
	"Events-Service/internal/http-server/handlers/event/createEvent"
	"Events-Service/internal/http-server/handlers/event/deleteEvent"
	"Events-Service/internal/http-server/handlers/event/getEvents"
	"Events-Service/internal/http-server/handlers/event/updateEvent"
	"Events-Service/internal/http-server/handlers/user"
	"Events-Service/internal/http-server/middleware/mwlogger"
	"Events-Service/internal/lib/logger/handlers/slogpretty"
	"Events-Service/internal/lib/logger/sl"
	"Events-Service/internal/storage/postgres"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("Starting events service", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	storage, err := postgres.InitDB(cfg)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(mwlogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Post("/create_user", user.New(log, storage))
	router.Post("/create_event", createEvent.New(log, storage))
	router.Post("/update_event", updateEvent.New(log, storage))
	router.Post("/delete_event", deleteEvent.New(log, storage))
	router.Get("/events_for_day", getEvents.ByDay(log, storage))
	router.Get("/events_for_week", getEvents.ByWeek(log, storage))
	router.Get("/events_for_month", getEvents.ByMonth(log, storage))

	log.Info("starting server", slog.String("address", cfg.HTTPServer.Address))

	srv := &http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err = srv.ListenAndServe(); err != nil {
		log.Error("failed to start server", sl.Err(err))
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sign := <-stop

	log.Info("server stopped", slog.String("signal", sign.String()))

	if err = storage.Close(); err != nil {
		log.Error("failed to close database", slog.String("error", err.Error()))
	}

	log.Info("postgres connection closed")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupPrettySlog()
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
