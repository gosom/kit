package es

import (
	"context"
	"time"

	"github.com/gosom/kit/logging"
)

type Worker interface {
	Process(ctx context.Context, key, value []byte, timestamp time.Time) error
}

type saveCommandWorker struct {
	log   logging.Logger
	store EventStore
	reg   *Registry
}

func NewSaveCommandWorker(store EventStore, registry *Registry) *saveCommandWorker {
	return &saveCommandWorker{
		log:   logging.Get().With("component", "save_command_worker"),
		store: store,
		reg:   registry,
	}
}

func (o *saveCommandWorker) Process(ctx context.Context, key, value []byte, timestamp time.Time) error {
	o.log.Info("Processing message", "key", string(key), "value", string(value), "timestamp", timestamp)
	busMsg := BusMessage{
		Key:       key,
		Data:      value,
		Timestamp: timestamp,
	}
	var cr CommandRecord
	if err := BusMessageToCommandRecord(busMsg, &cr); err != nil {
		return err
	}
	_, err := o.store.SaveCommandRecords(ctx, cr)
	return err
}
