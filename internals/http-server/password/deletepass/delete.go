package deletepass

import (
	"log/slog"
	"net/http"
	"password-db/internals/http-server/user/deleteuser"
	"password-db/internals/lib/api/response"
	"password-db/internals/storage/postgres"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type Request struct {
	ServiceName string `json:"service_name"`
}

func New(log *slog.Logger, s *postgres.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "internals.http-server.password.deletepass.New"

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

		rowsAffected, err := s.Delete(user, req.ServiceName)
		if err != nil {
			log.Error("failed to delete user from storage")

			render.JSON(w, r, response.Error("internal error"))

			return
		}

		render.JSON(w, r, deleteuser.DelResponse{
			Response:     response.OK(),
			RowsAffected: rowsAffected,
		})
	}
}
