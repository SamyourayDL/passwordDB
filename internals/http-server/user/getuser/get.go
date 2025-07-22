package getuser

import (
	"log/slog"
	"net/http"
	"password-db/internals/lib/api/response"
	resp "password-db/internals/lib/api/response"
	"password-db/internals/storage/postgres"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type GetResponse struct {
	resp.Response
	Passwords []postgres.PassResponse `json:"passwords"`
}

func New(log *slog.Logger, s *postgres.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "internals.http-server.user.getuser.New"

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

		passwords, err := s.GetPass(user, "")
		if err != nil {
			log.Error("failed to get user's passwords from storage")

			render.JSON(w, r, response.Error("internal error"))

			return
		}

		render.JSON(w, r, GetResponse{
			Response:  resp.OK(),
			Passwords: passwords,
		})
	}
}
