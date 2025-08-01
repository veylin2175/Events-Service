package createEvent

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
	UserId int64  `json:"user_id"`
	Date   string `json:"date"`
	Text   string `json:"text"`
}

type Response struct {
	response.Response
	EventId int64 `json:"event_id"`
}

type CreateEvent interface {
	SaveEvent(userID int64, dateStr, text string) (int64, error)
}

func New(log *slog.Logger, event CreateEvent) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.event.createEvent.New"

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

		eventId, err := event.SaveEvent(req.UserId, req.Date, req.Text)
		if errors.Is(err, storage.ErrEventExists) {
			log.Info("event already exists", slog.Int64("event", eventId))

			render.JSON(w, r, response.Error("event already exists"))

			return
		}
		if err != nil {
			log.Error("failed to add event", sl.Err(err))

			render.JSON(w, r, response.Error("failed to add event"))

			return
		}

		log.Info("event added", slog.Int64("id", eventId))

		responseOK(w, r, eventId)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, eventId int64) {
	render.JSON(w, r, Response{
		Response: response.OK(),
		EventId:  eventId,
	})
}
