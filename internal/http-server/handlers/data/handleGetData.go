package data

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/paniccaaa/test-task-biocad/internal/storage/postgres"
)

type Response struct {
	ID    string `json:"id"`
	Page  string `json:"page"`
	Limit string `json:"limit"`
}

func HandleGetData(log *slog.Logger, storage *postgres.PostgresStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.data.handleGetData"

		unitGUID := chi.URLParam(r, "unit_guid")
		page := r.URL.Query().Get("page")
		limit := r.URL.Query().Get("limit")

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		_ = log

		render.JSON(w, r, &Response{ID: unitGUID, Page: page, Limit: limit})
	}
}
