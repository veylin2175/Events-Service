package updateEvent

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
	UserId  int64  `json:"user_id"`
	EventId int64  `json:"event_id"`
	Date    string `json:"date"`
	Text    string `json:"text"`
}

type Response struct {
	response.Response
}

type UpdateEvent interface {
	UpdateEvent(userID, eventID int64, dateStr, text string) error
}

func New(log *slog.Logger, event UpdateEvent) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.event.updateEvent.New"

		log = log.With(
			slog.String("op", op),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			render.JSON(w, r, response.Error("failed to decode request"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err = validator.New().Struct(req); err != nil {
			var validateErr validator.ValidationErrors
			errors.As(err, &validateErr)

			log.Error("invalid request", sl.Err(err))

			render.JSON(w, r, response.ValidationError(validateErr))

			return
		}

		eventId := req.EventId
		err = event.UpdateEvent(req.UserId, req.EventId, req.Date, req.Text)
		if errors.Is(err, storage.ErrEventNotFound) {
			log.Info("event not found", slog.Int64("event", eventId))

			render.JSON(w, r, response.Error("event not found"))

			return
		}
		if err != nil {
			log.Error("failed to update event", sl.Err(err))

			render.JSON(w, r, response.Error("failed to update event"))

			return
		}

		log.Info("event updated", slog.Int64("id", eventId))

		responseOK(w, r)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, Response{
		Response: response.OK(),
	})
}
