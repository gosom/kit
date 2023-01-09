package todo

import (
	"context"
	"fmt"

	"github.com/gosom/kit/es"
)

var _ es.ICommand = (*CreateTodo)(nil)

type CreateTodo struct {
	es.CommandBase

	ID    string `json:"id" aggregateID:"true" validate:"required,uuid"`
	Title string `json:"title" validate:"required,gte=1,lte=140"`
}

func (c *CreateTodo) Handle(ctx context.Context, h es.AggregateLoader) ([]es.IEvent, error) {
	ev := TodoCreated{
		ID:    c.ID,
		Title: c.Title,
	}
	return []es.IEvent{&ev}, nil
}

type TodoCreated struct {
	es.EventBase

	ID    string `json:"id" aggregateID:"true"`
	Title string `json:"title"`
}

func (o *TodoCreated) Apply(aggregate es.AggregateRoot) error {
	agg, ok := aggregate.(*TodoAggregate)
	if !ok {
		return fmt.Errorf("%w but is %T", es.ErrInvalidAggregate, aggregate)
	}
	agg.SetID(o.ID)
	agg.Todo = NewTodo(o.ID)
	agg.Todo.ID = o.ID
	agg.Todo.Title = o.Title
	agg.Todo.Status = "open"
	return nil
}

type UpdateTodoStatus struct {
	es.CommandBase

	ID     string `json:"id" aggregateID:"true" validate:"required,uuid"`
	Status string `json:"status" validate:"required,oneof=open completed"`
}

func (c *UpdateTodoStatus) Handle(ctx context.Context, h es.AggregateLoader) ([]es.IEvent, error) {
	agg, err := LoadTodoAggregate(ctx, h, c.GetAggregateID())
	if err != nil {
		return nil, err
	}

	if err := agg.Todo.UpdateStatus(c.Status); err != nil {
		return nil, err
	}

	ev := TodoStatusUpdated{
		ID:     c.ID,
		Status: c.Status,
	}
	return []es.IEvent{&ev}, nil
}

type TodoStatusUpdated struct {
	es.EventBase

	ID     string `json:"id" aggregateID:"true"`
	Status string `json:"status"`
}

func (o *TodoStatusUpdated) Apply(aggregate es.AggregateRoot) error {
	agg, ok := aggregate.(*TodoAggregate)
	if !ok {
		return fmt.Errorf("%w but is %T", es.ErrInvalidAggregate, aggregate)
	}
	return agg.Todo.UpdateStatus(o.Status)
}
