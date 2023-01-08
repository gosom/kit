package mock

import (
	"context"

	"github.com/gosom/kit/es"
)

var _ es.EventStore = (*EventStore)(nil)

type EventStore struct {
}

func NewEventStore() *EventStore {
	return &EventStore{}
}

func (e *EventStore) GetCommand(ctx context.Context, commandID string) (es.CommandRecord, error) {
	panic("not implemented") // TODO: Implement
}

func (s *EventStore) Migrate(ctx context.Context) error {
	return nil
}

func (s *EventStore) SaveCommandRecords(ctx context.Context, records ...es.CommandRecord) ([]string, error) {
	return nil, nil
}

func (s *EventStore) SelectForProcessing(ctx context.Context, batchSize, max int) ([][]es.CommandRecord, error) {
	return nil, nil
}

func (s *EventStore) StoreCommandResults(ctx context.Context, commandID string,
	expectedVersion int, records ...es.EventRecord) error {
	return nil
}

func (s *EventStore) GetOrCreateVersion(ctx context.Context, aggregateID string) (int, error) {
	return 0, nil
}

func (s *EventStore) InsertSubscription(ctx context.Context, subscription string) (es.Subscription, error) {
	return es.Subscription{}, nil
}

func (s *EventStore) SelectEventsForSubscription(ctx context.Context, subscription es.Subscription, limit int) ([]es.EventRecord, error) {
	return nil, nil
}

func (s *EventStore) UpdateSubscription(ctx context.Context, group string, lastSeen string) (es.Subscription, error) {
	return es.Subscription{}, nil
}

func (s *EventStore) LoadEvents(ctx context.Context, aggregateID string) ([]es.EventRecord, error) {
	return nil, nil
}
