package user

import (
	"Events-Service/internal/lib/api/response"
	"Events-Service/internal/lib/logger/sl"
	"errors"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
)

type Request struct{}

type Response struct {
	response.Response
	UserId int64 `json:"user_id"`
}

type UserCreator interface {
	CreateUser() (int64, error)
}

func New(log *slog.Logger, userCreator UserCreator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.user.New"

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

		userId, err := userCreator.CreateUser()
		if err != nil {
			log.Error("failed to create user", sl.Err(err))

			render.JSON(w, r, response.Error("failed to create user"))

			return
		}

		log.Info("user created", slog.Int64("id", userId))

		responseOK(w, r, userId)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, id int64) {
	render.JSON(w, r, Response{
		Response: response.OK(),
		UserId:   id,
	})
}
