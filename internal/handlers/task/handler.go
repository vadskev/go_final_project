package task

import (
	"context"
	"encoding/json"
	"github.com/vadskev/go_final_project/internal/api"
	"github.com/vadskev/go_final_project/internal/logger"
	"github.com/vadskev/go_final_project/internal/nextdate"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/render"
	"github.com/pkg/errors"
	"github.com/vadskev/go_final_project/internal/models/task"
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
	const op = "transport.handlers.task.HandleGet"
	id := r.URL.Query().Get("id")

	if id == "" {
		api.ResponseError(w, r, errors.New("No id task, error ").Error(), http.StatusBadRequest)
		logger.Error(op, zap.Any("No id task:", "error"))
		return
	}

	dbTask, err := h.taskRepository.GetById(id)
	if err != nil {
		api.ResponseError(w, r, errors.New("No task, error ").Error(), http.StatusNotFound)
		logger.Error(op, zap.Any("No task:", "error"))
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, dbTask)
}

func (h *Handler) HandlePost(w http.ResponseWriter, r *http.Request) {
	const op = "transport.handlers.task.HandlePost"

	var taskInfo task.Info

	err := render.DecodeJSON(r.Body, &taskInfo)
	if err != nil {
		api.ResponseError(w, r, err.Error(), http.StatusBadRequest)
		return
	}

	if taskInfo.Title == "" {
		api.ResponseError(w, r, errors.New("No title").Error(), http.StatusBadRequest)
		return
	}

	now := time.Now()
	startDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	if taskInfo.Date == "" {
		taskInfo.Date = now.Format("20060102")
	}

	date, err := time.Parse("20060102", taskInfo.Date)
	if err != nil {
		api.ResponseError(w, r, errors.New("No valid date format").Error(), http.StatusBadRequest)
		logger.Error(op, zap.Any("No valid date format, error:", err.Error()))
		return
	}

	if date.Before(startDay) {
		if taskInfo.Repeat == "" {
			taskInfo.Date = now.Format("20060102")
		} else {
			taskInfo.Date, err = nextdate.NextDate(now, taskInfo.Date, taskInfo.Repeat)
			if err != nil {
				api.ResponseError(w, r, errors.New("No correct repeat format").Error(), http.StatusBadRequest)
				logger.Error(op, zap.Any("No correct repeat format, error:", err.Error()))
				return
			}
		}
	}

	var response task.Response
	response.ID, err = h.taskRepository.Create(&taskInfo)
	if err != nil {
		api.ResponseError(w, r, errors.New("Error creating task: ").Error()+err.Error(), http.StatusBadRequest)
		logger.Error(op, zap.Any("Error creating task, error:", err.Error()))
		return
	}
	api.ResponseOK(w, r, response)
}

func (h *Handler) HandlePut(w http.ResponseWriter, r *http.Request) {
	const op = "transport.handlers.task.HandlePut"

	var newTask task.Task

	err := json.NewDecoder(r.Body).Decode(&newTask)
	if err != nil {
		api.ResponseError(w, r, errors.New("Deserialize json").Error(), http.StatusBadRequest)
		logger.Error(op, zap.Any("Deserialize json, error:", err.Error()))
		return
	}

	if newTask.ID == "" {
		api.ResponseError(w, r, errors.New("No task id").Error(), http.StatusBadRequest)
		logger.Error(op, zap.Any("No task id, error:", errors.New("No task id").Error()))
		return
	}

	_, err = strconv.Atoi(newTask.ID)
	if err != nil {
		api.ResponseError(w, r, errors.New("ID not number").Error(), http.StatusBadRequest)
		logger.Error(op, zap.Any("ID not number, error:", err.Error()))
		return
	}

	if newTask.Title == "" {
		api.ResponseError(w, r, errors.New("Title empty").Error(), http.StatusBadRequest)
		logger.Error(op, zap.Any("Title empty, error:", errors.New("Title empty").Error()))
		return
	}

	now := time.Now()
	startDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	if newTask.Date == "" {
		newTask.Date = now.Format("20060102")
	}

	date, err := time.Parse("20060102", newTask.Date)
	if err != nil {
		api.ResponseError(w, r, errors.New("No valid date format").Error(), http.StatusBadRequest)
		logger.Error(op, zap.Any("No valid date format, error:", err.Error()))
		return
	}

	if date.Before(startDay) {
		if newTask.Repeat == "" {
			newTask.Date = now.Format("20060102")
		} else {
			newTask.Date, err = nextdate.NextDate(now, newTask.Date, newTask.Repeat)
			if err != nil {
				api.ResponseError(w, r, errors.New("No correct repeat format").Error(), http.StatusBadRequest)
				logger.Error(op, zap.Any("No correct repeat format, error:", err.Error()))
				return
			}
		}
	}

	err = h.taskRepository.Update(&newTask)
	if err != nil {
		api.ResponseError(w, r, errors.New("Error update Task").Error(), http.StatusInternalServerError)
		logger.Error(op, zap.Any("Error update Task, error:", err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	api.ResponseOK(w, r, task.Response{})
}

func (h *Handler) HandleDelete(w http.ResponseWriter, r *http.Request) {
	const op = "transport.handlers.task.HandlePut"
	id := r.URL.Query().Get("id")
	if id == "" {
		api.ResponseError(w, r, errors.New("No id task, error ").Error(), http.StatusBadRequest)
		logger.Error(op, zap.Any("No id task:", "error"))
		return
	}

	err := h.taskRepository.Delete(id)
	if err != nil {
		api.ResponseError(w, r, errors.New("Error delete task, error ").Error(), http.StatusNotFound)
		logger.Error(op, zap.Any("Error delete task:", "error"))
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, task.Response{})
}
