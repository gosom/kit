package postgres

import (
	"context"
	"time"

	"github.com/gosom/kit/es"
	"github.com/gosom/kit/logging"
)

type worker struct {
	log   logging.Logger
	store es.EventStore
}

func NewWorker(store es.EventStore) *worker {
	return &worker{
		log:   logging.Get().With("component", "worker"),
		store: store,
	}
}

func (w *worker) Process(ctx context.Context, key []byte, value []byte, timestamp time.Time) error {
	w.log.Info("Processing message", "key", string(key), "value", string(value), "timestamp", timestamp)
	return nil
}
