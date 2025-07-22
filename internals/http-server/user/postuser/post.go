package postuser

import (
	"log/slog"
	"net/http"
	"password-db/internals/lib/api/response"
	"password-db/internals/storage/postgres"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

func New(log *slog.Logger, s *postgres.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "internals.http-server.user.post.New"

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
		if err := s.AddUser(user); err != nil {
			log.Error("failed to add user to storage", "err", err.Error())

			render.JSON(w, r, response.Error("failed to add user to storage"))

			return
		}

		log.Info("user was successfully aded", "user_name", user)

		render.JSON(w, r, response.OK())
	}
}
