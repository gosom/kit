package es

import "context"

// EventStore is the interface that wraps the basic event store methods.
type EventStore interface {
	//Migrate runs the migrations for the event store.
	Migrate(ctx context.Context) error

	//SaveCommandRecord saves the command record.
	SaveCommandRecords(ctx context.Context, records ...CommandRecord) ([]string, error)
	GetCommand(ctx context.Context, commandID string) (CommandRecord, error)

	//StoreCommandResults stores the command results.
	StoreCommandResults(ctx context.Context, commandID string, expectedVersion int, events ...EventRecord) error

	//SelectForProcessing selects the command records for processing.
	SelectForProcessing(ctx context.Context, workers, limit int) ([][]CommandRecord, error)

	//GetOrCreateVersion gets the version for the aggregate or creates it if it doesn't exist.
	GetOrCreateVersion(ctx context.Context, aggregateID string) (int, error)

	//InsertSubscription
	InsertSubscription(ctx context.Context, subscription string) (Subscription, error)
	//SelectEventsForSubscription selects the events for the subscription.
	SelectEventsForSubscription(ctx context.Context, subscription Subscription, limit int) ([]EventRecord, error)
	//UpdateSubscription updates the subscription.
	UpdateSubscription(ctx context.Context, group string, lastSeen string) (Subscription, error)

	//LoadEvents loads the events for the aggregate.
	LoadEvents(ctx context.Context, aggregateID string) ([]EventRecord, error)
}
