package signin

import (
	"context"
	"github.com/go-chi/render"
	"github.com/pkg/errors"
	"github.com/vadskev/go_final_project/internal/api"
	"github.com/vadskev/go_final_project/internal/config"
	"github.com/vadskev/go_final_project/internal/logger"
	"github.com/vadskev/go_final_project/internal/storage/db"
	"go.uber.org/zap"
	"net/http"
)

type Handler struct {
	taskRepository db.Repository
	ctx            context.Context
	pass           config.PasswordConfig
}

func New(ctx context.Context, taskRepository db.Repository, pass config.PasswordConfig) *Handler {
	return &Handler{
		taskRepository: taskRepository,
		ctx:            ctx,
		pass:           pass,
	}
}

type SignRequest struct {
	Password string `json:"password"`
}

type SignResponse struct {
	Token string `json:"token,omitempty"`
	Error string `json:"error,omitempty"`
}

func (h *Handler) HandlePost(w http.ResponseWriter, r *http.Request) {
	const op = "transport.handlers.signin.HandlePost"

	var req SignRequest
	err := render.DecodeJSON(r.Body, &req)
	if err != nil {
		api.ResponseError(w, r, err.Error(), http.StatusBadRequest)
		logger.Error(op, zap.Any("Decode json, error:", err.Error()))
		return
	}

	if req.Password != h.pass.GetPass() {
		api.ResponseError(w, r, errors.New("Wrong password").Error(), http.StatusBadRequest)
		logger.Error(op, zap.Any("error:", errors.New("Wrong password").Error()))
		return
	}

	tokenString := h.pass.CreateHash(req.Password)
	resp := SignResponse{Token: tokenString}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, resp)
}
