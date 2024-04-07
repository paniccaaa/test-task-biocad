package data

import (
	"log/slog"
	"net/http"
	"strconv"

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

type ErrorResponse struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

func HandleGetDataByID(log *slog.Logger, storage *postgres.PostgresStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.data.handleGetDataByID"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		unitGUID := chi.URLParam(r, "unit_guid")
		page, err := strconv.Atoi(r.URL.Query().Get("page"))
		if err != nil {
			log.Error("failed to convert page to int", slog.String("err", err.Error()))
		}
		limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
		if err != nil {
			log.Error("failed to convert limit to int", slog.String("err", err.Error()))
		}

		offset := (page - 1) * limit

		data, err := storage.GetDataByUnitGUIDWithPagination(unitGUID, limit, offset)
		if err != nil {
			log.Error("failed to get data with pagination", slog.String("err", err.Error()))

			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, ErrorResponse{Status: "Error", Error: err.Error()})

			return
		}

		render.JSON(w, r, data)
	}
}
