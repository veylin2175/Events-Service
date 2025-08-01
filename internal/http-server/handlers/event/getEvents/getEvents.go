package getEvents

import (
	"Events-Service/internal/lib/api/response"
	"Events-Service/internal/lib/logger/sl"
	"Events-Service/internal/models"
	"errors"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
	"time"
)

type EventResponse struct {
	Date string `json:"date"`
	Text string `json:"text"`
}

type Request struct {
	UserId int64  `json:"user_id"`
	Date   string `json:"date"`
}

type Response struct {
	response.Response
	Events []EventResponse `json:"events"`
}

type GetEvents interface {
	GetEventsByDay(userID int64, date string) ([]models.Event, error)
	GetEventsByWeek(userID int64, date time.Time) ([]models.Event, error)
	GetEventsByMonth(userID int64, year int, month time.Month) ([]models.Event, error)
}

func ByDay(log *slog.Logger, event GetEvents) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.event.getEvents.New"

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

		events, err := event.GetEventsByDay(req.UserId, req.Date)
		if err != nil {
			log.Error("failed to get events", sl.Err(err))

			render.JSON(w, r, response.Error("failed to get events"))

			return
		}

		responseEvents := make([]EventResponse, 0, len(events))
		for _, e := range events {
			responseEvents = append(responseEvents, EventResponse{
				Date: e.Date,
				Text: e.Text,
			})
		}

		log.Info("got events")

		responseOK(w, r, responseEvents)
	}
}

func ByWeek(log *slog.Logger, event GetEvents) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.event.getEvents.ByWeek"

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

		parsedDate, err := time.Parse("2006-01-02", req.Date)
		if err != nil {
			log.Error("invalid date format", sl.Err(err))
			render.JSON(w, r, response.Error("invalid date format, use YYYY-MM-DD"))
			return
		}

		events, err := event.GetEventsByWeek(req.UserId, parsedDate)
		if err != nil {
			log.Error("failed to get events", sl.Err(err))

			render.JSON(w, r, response.Error("failed to get events"))

			return
		}

		responseEvents := make([]EventResponse, 0, len(events))
		for _, e := range events {
			responseEvents = append(responseEvents, EventResponse{
				Date: e.Date,
				Text: e.Text,
			})
		}

		log.Info("got events")

		responseOK(w, r, responseEvents)
	}
}

func ByMonth(log *slog.Logger, event GetEvents) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.event.getEvents.ByMonth"

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

		parsedDate, err := time.Parse("2006-01-02", req.Date)
		if err != nil {
			log.Error("invalid date format", sl.Err(err))
			render.JSON(w, r, response.Error("invalid date format, use YYYY-MM-DD"))
			return
		}

		year := parsedDate.Year()
		month := parsedDate.Month()

		events, err := event.GetEventsByMonth(req.UserId, year, month)
		if err != nil {
			log.Error("failed to get events", sl.Err(err))

			render.JSON(w, r, response.Error("failed to get events"))

			return
		}

		responseEvents := make([]EventResponse, 0, len(events))
		for _, e := range events {
			responseEvents = append(responseEvents, EventResponse{
				Date: e.Date,
				Text: e.Text,
			})
		}

		log.Info("got events")

		responseOK(w, r, responseEvents)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, events []EventResponse) {
	render.JSON(w, r, Response{
		Response: response.OK(),
		Events:   events,
	})
}
