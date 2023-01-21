package es

import (
	"context"
	"fmt"
)

type AggregateLoader interface {
	Load(ctx context.Context, aggregateID string, agg AggregateRoot) error
}

type AggregateFactory func() (AggregateRoot, error)

type AggregateRoot interface {
	GetID() string
	SetID(string)

	GetType() string
	SetType(string)

	GetVersion() uint64

	SetVersion(uint64)

	String() string
}

type AggregateBase struct {
	ID      string
	Type    string
	Version uint64
}

func NewAggregateBase() (*AggregateBase, error) {
	return &AggregateBase{
		Version: 0,
	}, nil
}

func (a *AggregateBase) GetID() string {
	return a.ID
}

func (a *AggregateBase) SetID(id string) {
	a.ID = id
}

func (a *AggregateBase) GetType() string {
	return a.Type
}

func (a *AggregateBase) SetType(t string) {
	a.Type = t
}

func (a *AggregateBase) GetVersion() uint64 {
	return a.Version
}

func (a *AggregateBase) String() string {
	return fmt.Sprintf("AggregateBase{ID: %s, Type: %s, Version: %d}", a.ID, a.Type, a.Version)
}

func (a *AggregateBase) SetVersion(v uint64) {
	a.Version = v
}

func Load(agg AggregateRoot, events []IEvent) error {
	for i := range events {
		switch e := events[i].(type) {
		case *EventError:
		default:
			if err := RaiseEvent(agg, e); err != nil {
				return err
			}
		}
	}
	return nil
}

func RaiseEvent(agg AggregateRoot, event IEvent) error {
	if agg == nil {
		return ErrNilAggregate
	}
	if event == nil {
		return ErrInvalidEvent
	}
	if int(agg.GetVersion()) >= event.GetVersion() {
		return fmt.Errorf("event version %d is not greater than aggregate version %d: %w", event.GetVersion(), agg.GetVersion(), ErrDuplicateEvent)
	}
	if err := event.Apply(agg); err != nil {
		return err
	}
	agg.SetVersion(uint64(event.GetVersion()))
	return nil
}
