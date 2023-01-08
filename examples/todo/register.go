package todo

import (
	"github.com/gosom/kit/es"
)

func Register(registry *es.Registry) {
	serializer := Serializer{}
	registry.RegisterCommand("CreateTodo", serializer.CreateTodo)
	registry.RegisterEvent("TodoCreated", serializer.TodoCreated)

	registry.RegisterCommand("UpdateTodoStatus", serializer.UpdateTodoStatus)
	registry.RegisterEvent("TodoStatusUpdated", serializer.TodoStatusUpdated)
}
