package es

import "context"

// Publisher is an interface for publishing events.
type Publisher interface {
	Name() string
	Publish(ctx context.Context, events ...EventRecord) error
}
