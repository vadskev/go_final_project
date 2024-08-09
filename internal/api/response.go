package api

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/vadskev/go_final_project/internal/models/task"
)

type Response struct {
	Status int    `json:"status,omitempty"`
	Error  string `json:"error,omitempty"`
}

func ResponseError(w http.ResponseWriter, r *http.Request, error string, status int) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(status)
	render.JSON(w, r, Response{
		Error: error,
	})
}
func ResponseOK(w http.ResponseWriter, r *http.Request, response task.Response) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, response)
}
