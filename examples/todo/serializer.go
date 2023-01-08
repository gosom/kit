package todo

import (
	"encoding/json"

	"github.com/gosom/kit/es"
)

type Serializer struct {
}

func (s *Serializer) CreateTodo(rec es.CommandRecord) (es.ICommand, error) {
	var item CreateTodo
	err := json.Unmarshal(rec.Data, &item)
	return &item, err
}

func (s *Serializer) TodoCreated(rec es.EventRecord) (es.IEvent, error) {
	var item TodoCreated
	err := json.Unmarshal(rec.Data, &item)
	return &item, err
}

func (s *Serializer) UpdateTodoStatus(rec es.CommandRecord) (es.ICommand, error) {
	var item UpdateTodoStatus
	err := json.Unmarshal(rec.Data, &item)
	return &item, err
}

func (s *Serializer) TodoStatusUpdated(rec es.EventRecord) (es.IEvent, error) {
	var item TodoStatusUpdated
	err := json.Unmarshal(rec.Data, &item)
	return &item, err
}
