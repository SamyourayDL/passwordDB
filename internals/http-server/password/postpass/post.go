package postpass

import (
	"log/slog"
	"net/http"
	"password-db/internals/lib/api/response"
	"password-db/internals/storage/postgres"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type Request struct {
	Password    string `json:"password"`
	ServiceName string `json:"service_name"`
	Category    string `json:"omitempty"`
}

func New(log *slog.Logger, s *postgres.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "internals.http-server.password.postpass.New"

		log := log.With(
			"fn", fn,
			"request_id", middleware.GetReqID(r.Context()),
		)

		user := chi.URLParam(r, "user_name")
		if user == "" {
			log.Error("failed to get user_name from url")

			render.JSON(w, r, response.Error("internal error"))

			return
		}

		var req Request

		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("failed to parse request body")

			render.JSON(w, r, response.Error("internal error"))

			return
		}

		if err := s.AddPassword(user, req.Password, req.ServiceName, req.Category); err != nil {
			log.Error("failed to add password to storage", "err", err.Error())

			render.JSON(w, r, response.Error("failed to add password to storage"))

			return
		}

		render.JSON(w, r, response.OK())
	}
}
