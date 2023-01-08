package todo

import (
	"context"

	"github.com/gosom/kit/es"
)

type TodoAggregate struct {
	*es.AggregateBase
	Todo Todo
}

func NewTodoAggregate() (es.AggregateRoot, error) {
	agg := TodoAggregate{}
	base, err := es.NewAggregateBase()
	if err != nil {
		return nil, err
	}
	base.SetType("todo")
	agg.AggregateBase = base
	return &agg, nil
}

func LoadTodoAggregate(ctx context.Context, h es.AggregateLoader, id string) (*TodoAggregate, error) {
	agg, err := NewTodoAggregate()
	if err != nil {
		return nil, err
	}
	if err := h.Load(ctx, id, agg); err != nil {
		return nil, err
	}
	return agg.(*TodoAggregate), nil
}
