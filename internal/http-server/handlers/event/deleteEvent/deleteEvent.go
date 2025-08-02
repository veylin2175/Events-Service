package deleteEvent

import (
	"Events-Service/internal/lib/api/response"
	"Events-Service/internal/lib/logger/sl"
	"Events-Service/internal/storage"
	"errors"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
)

type Request struct {
	UserId  int64 `json:"user_id" validate:"required"`
	EventId int64 `json:"event_id" validate:"required"`
}

type Response struct {
	response.Response
}

//go:generate go run github.com/vektra/mockery/v2@v2.51.1 --name=DeleteEvent
type DeleteEvent interface {
	DeleteEvent(userID, eventID int64) error
}

func New(log *slog.Logger, event DeleteEvent) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.event.deleteEvent.New"

		log = log.With(
			slog.String("op", op),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("failed to decode request"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err = validator.New().Struct(req); err != nil {
			var validateErr validator.ValidationErrors
			errors.As(err, &validateErr)

			log.Error("invalid request", sl.Err(err))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.ValidationError(validateErr))

			return
		}

		eventId := req.EventId
		err = event.DeleteEvent(req.UserId, req.EventId)
		if errors.Is(err, storage.ErrEventNotFound) {
			log.Info("event not found", slog.Int64("event", eventId))
			render.Status(r, http.StatusServiceUnavailable)
			render.JSON(w, r, response.Error("event not found"))

			return
		}
		if err != nil {
			log.Error("failed to delete event", sl.Err(err))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed to delete event"))

			return
		}

		log.Info("event deleted", slog.Int64("id", eventId))

		responseOK(w, r)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, Response{
		Response: response.OK(),
	})
}
