package es

import (
	"errors"
)

var (
	ErrInvalidCommand = errors.New("invalid command")

	ErrUnknownEventStoreType       = errors.New("unknown event store type")
	ErrUnknownCommandListenerType  = errors.New("unknown command listener type")
	ErrUnknownCommandProcessorType = errors.New("unknown command processor type")
	ErrWrongExpectedVersion        = errors.New("wrong expected version")
	ErrUnknownPublisherType        = errors.New("unknown publisher type")
	ErrUnknownSubscriberType       = errors.New("unknown subscriber type")

	ErrInvalidEvent      = errors.New("invalid event")
	ErrUnregisteredEvent = errors.New("unregistered event")

	ErrSkipEvent        = errors.New("skip event")
	ErrInvalidAggregate = errors.New("invalid aggregate")

	ErrNilAggregate = errors.New("nil aggregate")
)

type EventError struct {
	EventBase

	Err error
}

func NewEventErrorFromError(err error, expectedVersion int, rec CommandRecord) EventError {
	ev := EventError{
		Err: err,
	}
	return ev
}
