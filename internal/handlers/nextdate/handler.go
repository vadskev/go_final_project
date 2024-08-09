package nextdate

import (
	"context"
	"github.com/vadskev/go_final_project/internal/api"
	"github.com/vadskev/go_final_project/internal/logger"
	"github.com/vadskev/go_final_project/internal/nextdate"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/vadskev/go_final_project/internal/storage/db"
	"go.uber.org/zap"
)

type Handler struct {
	taskRepository db.Repository
	ctx            context.Context
}

func New(ctx context.Context, taskRepository db.Repository) *Handler {
	return &Handler{
		taskRepository: taskRepository,
		ctx:            ctx,
	}
}

func (h *Handler) HandleGet(w http.ResponseWriter, r *http.Request) {
	const op = "transport.handlers.nextdate.Handle"
	nowParam := r.URL.Query().Get("now")
	dateParam := r.URL.Query().Get("date")
	repeatParam := r.URL.Query().Get("repeat")

	now, err := time.Parse("20060102", nowParam)
	if err != nil {
		api.ResponseError(w, r, errors.New("Vrong NOW param").Error(), http.StatusBadRequest)
		logger.Error(op, zap.Any("Vrong NOW param, error:", errors.New("No taskItem id").Error()))
		return
	}

	nextDate, err := nextdate.NextDate(now, dateParam, repeatParam)
	if err != nil {
		api.ResponseError(w, r, errors.New("No find next date").Error(), http.StatusInternalServerError)
		logger.Error(op, zap.Any("No find next date, error:", errors.New("No taskItem id").Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(nextDate))
	if err != nil {
		logger.Error(op, zap.Error(err))
	}
}
