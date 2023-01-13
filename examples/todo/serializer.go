package todo

import (
	"encoding/json"

	"github.com/gosom/kit/es"
)

type Serializer struct {
}

func (s *Serializer) CreateTodo(data []byte) (es.ICommand, error) {
	var item CreateTodo
	return &item, json.Unmarshal(data, &item)
}

func (s *Serializer) TodoCreated(data []byte) (es.IEvent, error) {
	var item TodoCreated
	return &item, json.Unmarshal(data, &item)
}

func (s *Serializer) UpdateTodoStatus(data []byte) (es.ICommand, error) {
	var item UpdateTodoStatus
	return &item, json.Unmarshal(data, &item)
}

func (s *Serializer) TodoStatusUpdated(data []byte) (es.IEvent, error) {
	var item TodoStatusUpdated
	return &item, json.Unmarshal(data, &item)
}
