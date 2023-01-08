package api

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gosom/kit/core"
	"github.com/gosom/kit/es"
	"github.com/gosom/kit/examples/todo"
	"github.com/gosom/kit/web"
)

func RegisterHandlers(r web.Router, dispatcher es.CommandDispatcher) {
	handler := NewTodoHandler(dispatcher)
	r.Post("/todo", handler.CreateTodo)
	r.Patch("/todo/{id}", handler.UpdateTodoStatus)
}

type TodoHandler struct {
	dispatcher es.CommandDispatcher
}

func NewTodoHandler(dispatcher es.CommandDispatcher) *TodoHandler {
	return &TodoHandler{dispatcher: dispatcher}
}

// CreateTodoRequest is the request payload for the CreateTodo method.
type CreateTodoRequest struct {
	Title string `json:"title"`
}

// CreateTodoResponse is the response payload for the CreateTodo method.
type CreateTodoResponse struct {
	CommandID string `json:"commandId"`
	ID        string `json:"id"`
}

func (o *TodoHandler) CreateTodo(w http.ResponseWriter, r *http.Request) {
	var payload CreateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		web.JSONError(w, r, core.ErrBadRequest)
		return
	}
	var ans CreateTodoResponse
	var err error
	cmd := todo.CreateTodo{
		ID:    uuid.New().String(),
		Title: payload.Title,
	}
	ans.ID = cmd.ID
	ans.CommandID, err = o.dispatcher.Dispatch(r.Context(), &cmd)
	if err != nil {
		web.JSONError(w, r, err)
		return
	}
	web.JSON(w, r, http.StatusAccepted, ans)
}

// UpdateTodoStatusRequest is the request payload for the CompleteTodo method.
type UpdateTodoStatusRequest struct {
	ID     string `json:"-"`
	Status string `json:"status"`
}

// UpdateTodoStatusResponse is the response payload for the CompleteTodo method.
type UpdateTodoStatusResponse struct {
	CommandID string `json:"commandId"`
}

func (o *TodoHandler) UpdateTodoStatus(w http.ResponseWriter, r *http.Request) {
	id := web.StringURLParam(r, "id")
	if len(id) == 0 {
		web.JSONError(w, r, core.ErrBadRequest)
		return
	}
	var payload UpdateTodoStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		web.JSONError(w, r, core.ErrBadRequest)
		return
	}
	payload.ID = id
	var ans UpdateTodoStatusResponse
	var err error
	cmd := todo.UpdateTodoStatus{
		ID:     payload.ID,
		Status: payload.Status,
	}
	ans.CommandID, err = o.dispatcher.Dispatch(r.Context(), &cmd)
	if err != nil {
		web.JSONError(w, r, err)
		return
	}
	web.JSON(w, r, http.StatusAccepted, ans)
}
