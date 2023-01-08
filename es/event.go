package es

import (
	"encoding/json"
	"fmt"
)

// EventBase is the base struct for all events.
type EventBase struct {
	CommandEventBase
	version int
}

// validate validates the event.
func (e *EventBase) Validate() error {
	if len(e.id) == 0 {
		return fmt.Errorf("Event ID is required %w", ErrInvalidEvent)
	}
	if len(e.aggregateID) == 0 {
		return fmt.Errorf("Event AggregateID is required %w", ErrInvalidEvent)
	}
	if len(e.eventType) == 0 {
		return fmt.Errorf("Event Type is required %w", ErrInvalidEvent)
	}
	return nil
}

func (e *EventBase) GetVersion() int {
	return e.version
}

func (e *EventBase) SetVersion(version int) {
	e.version = version
}

// ApplyEvent applies the event to the aggregate.
func (e *EventBase) Apply(aggregate AggregateRoot) error {
	panic("not implemented")
}

type EventRecord struct {
	RecordBase
	CommandID string
	Version   int
}

func (o *EventRecord) Bind() []any {
	ans := o.RecordBase.Bind()
	return append(ans, &o.CommandID, &o.Version)
}

func EventToEventRecord(ev IEvent) (EventRecord, error) {
	data, err := json.Marshal(ev)
	if err != nil {
		return EventRecord{}, err
	}
	return EventRecord{
		RecordBase: RecordBase{
			ID:          ev.GetID(),
			AggregateID: ev.GetAggregateID(),
			EventType:   ev.GetEventType(),
			Data:        data,
		},
		Version: ev.GetVersion(),
	}, nil
}

func EventRecordsToEvents(registry *Registry, records []EventRecord) ([]IEvent, error) {
	var err error
	ans := make([]IEvent, len(records))
	for i := range records {
		ans[i], err = EventRecordToEvent(registry, records[i])
		if err != nil {
			return nil, err
		}
	}
	return ans, nil
}

func EventRecordToEvent(registry *Registry, record EventRecord) (IEvent, error) {
	convFn, ok := registry.GetEvent(record.EventType)
	if !ok {
		return nil, fmt.Errorf("event type %s not found in registry", record.EventType)
	}
	ev, err := convFn(record)
	if err != nil {
		return nil, err
	}
	ev.SetVersion(record.Version)
	ev.SetID(record.ID)
	ev.SetEventType(record.EventType)
	ev.SetVersion(record.Version)
	ev.SetAggregateID(record.AggregateID)
	return ev, nil
}
