package postgres

import (
	"context"

	"github.com/gosom/kit/es"
	"github.com/gosom/kit/logging"
)

var _ es.CommandDispatcher = (*commandDispatcher)(nil)

type commandDispatcher struct {
	log    logging.Logger
	domain string
	store  es.EventStore
}

func NewCommandDispatcher(domain string, store es.EventStore) *commandDispatcher {
	dispatcher := commandDispatcher{
		log:    logging.Get().With("component", "command_dispatcher").Level(logging.DEBUG),
		domain: domain,
		store:  store,
	}
	return &dispatcher
}

func (c *commandDispatcher) Dispatch(ctx context.Context, cmd es.ICommand) (string, error) {
	cr, err := es.CommandToCommandRecord(c.domain, cmd)
	if err != nil {
		return "", err
	}
	records, err := c.store.SaveCommandRecords(ctx, cr)
	if err != nil {
		return "", err
	}
	if len(records) == 1 {
		c.log.Debug("command dispatched", "command", cmd, "command_id", records[0])
	}
	return records[0], nil
}
