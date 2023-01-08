package api

import (
	"github.com/gosom/kit/web"
)

func RegisterHandlers(r web.Router) {
	//handler := NewTodoHandler(dispatcher)
}

type TodoHandler struct {
}

func NewTodoHandler() *TodoHandler {
	return &TodoHandler{}
}
