package es

import (
	"context"
	"time"
)

// ICommand is a common interface for commands and events
type ICommandEvent interface {
	GetID() string
	SetID(id string)
	GetAggregateID() string
	SetAggregateID(aggregateID string)
	GetEventType() string
	SetEventType(eventType string)
	Validate() error
}

// ICommand is the interface that all commands must implement.
type ICommand interface {
	ICommandEvent
	GetAggregateHash() int32
	SetAggregateHash()
	Handle(ctx context.Context, h AggregateLoader) ([]IEvent, error)
}

// IEvent is the interface that all events must implement.
type IEvent interface {
	ICommandEvent
	SetVersion(version int)
	GetVersion() int
	Apply(aggregate AggregateRoot) error
}

// RecordBase is the base struct for all records (command & events)
type RecordBase struct {
	ID          string
	AggregateID string
	EventType   string
	Data        []byte
	CreatedAt   time.Time
}

func (o *RecordBase) Bind() []any {
	return []any{
		&o.ID,
		&o.AggregateID,
		&o.EventType,
		&o.Data,
		&o.CreatedAt,
	}
}

var _ ICommandEvent = (*CommandEventBase)(nil)

// CommandEventBase is the base struct for all commands and events.
type CommandEventBase struct {
	id          string
	aggregateID string
	eventType   string
}

func (c *CommandEventBase) GetID() string {
	return c.id
}

func (c *CommandEventBase) SetID(id string) {
	c.id = id
}

func (c *CommandEventBase) GetAggregateID() string {
	return c.aggregateID
}

func (c *CommandEventBase) SetAggregateID(aggregateID string) {
	c.aggregateID = aggregateID
}

func (c *CommandEventBase) GetEventType() string {
	return c.eventType
}

func (c *CommandEventBase) SetEventType(eventType string) {
	c.eventType = eventType
}

func (c *CommandEventBase) Validate() error {
	panic("not implemented")
}
