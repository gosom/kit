package es

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/gosom/kit/lib"
	"github.com/gosom/kit/logging"
)

// CommandProcessor is a processor for commands
// It is responsible to process commands
type CommandProcessor interface {
	AggregateLoader
	Start(ctx context.Context) error
}

var _ CommandProcessor = (*commandProcessor)(nil)

type commandProcessor struct {
	store     EventStore
	reg       *Registry
	workerNum int
	domain    string
	log       logging.Logger
}

func NewCommandProcessor(
	workerNum int,
	store EventStore,
	registry *Registry,
	domain string,
) (CommandProcessor, error) {
	var ans commandProcessor
	ans.workerNum = workerNum
	ans.store = store
	ans.reg = registry
	ans.domain = domain
	ans.log = logging.Get().With("component", "command_processor")
	return &ans, nil
}

func (c *commandProcessor) Start(ctx context.Context) error {
	c.log.Info("starting command processor")
	defer c.log.Info("command processor stopped")
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			if err := c.work(ctx); err != nil {
				c.log.Error("failed to process commands", "error", err)
			}
		}
	}
}

func (c *commandProcessor) Load(ctx context.Context, aggregateID string, aggregate AggregateRoot) error {
	records, err := c.store.LoadEvents(ctx, aggregateID)
	if err != nil {
		return err
	}
	events := make([]IEvent, len(records))
	for i := range records {
		convFn, ok := c.reg.GetEvent(records[i].EventType)
		if !ok {
			return err
		}
		events[i], err = convFn(records[i])
	}
	if err := Load(aggregate, events); err != nil {
		return err
	}
	return nil
}

func (c *commandProcessor) work(ctx context.Context) error {
	t0 := time.Now()
	items, err := c.store.SelectForProcessing(ctx, c.workerNum, 10)
	if err != nil {
		return fmt.Errorf("%w when selecting commands", err)
	}
	t1 := time.Now()
	g, ctx := errgroup.WithContext(ctx)
	total := 0
	for i := 0; i < len(items); i++ {
		total += len(items[i])
		if len(items[i]) > 0 {
			num := i
			g.Go(func() error {
				return c.processGroup(ctx, items[num])
			})
		}
	}
	if err := g.Wait(); err != nil {
		return fmt.Errorf("%w when processing commands", err)
	}
	t2 := time.Now()
	speed := float64(total) / t2.Sub(t0).Seconds()
	if total > 0 {
		c.log.Debug(
			"processed commands",
			"total_duration", t2.Sub(t0),
			"select_duration", t1.Sub(t0),
			"process_duration", t2.Sub(t1),
			"speed", speed,
			"total", total,
		)
	}
	return nil
}

func (c *commandProcessor) processGroup(ctx context.Context, items []CommandRecord) error {
	for i := 0; i < len(items); i++ {
		select {
		case <-ctx.Done():
			return nil
		default:
		}
		if err := c.process(ctx, items[i]); err != nil {
			return err
		}
	}
	return nil
}

func (c *commandProcessor) process(ctx context.Context, rec CommandRecord) (err error) {
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			err, ok = r.(error)
			if !ok {
				err = errors.New("unknown error")
			}
		}
	}()
	convFn, ok := c.reg.GetCommand(rec.EventType)
	if !ok {
		err = fmt.Errorf("no converter for event type %s %w", rec.EventType, ErrSkipEvent)
		return
	}
	expectedVersion, err := c.store.GetOrCreateVersion(ctx, rec.AggregateID)
	if err != nil {
		return
	}
	var cmd ICommand
	cmd, err = convFn(rec)
	if err != nil {
		err = fmt.Errorf("rec: %s %w", rec.EventType, err)
		return
	}

	cmd.SetID(rec.ID)
	cmd.SetEventType(rec.EventType)
	cmd.SetAggregateID(rec.AggregateID)
	cmd.SetAggregateHash()

	var newEvents []IEvent
	newEvents, err = cmd.Handle(ctx, c)
	if err != nil {
		errorEvent := NewEventErrorFromError(err, expectedVersion, rec)
		newEvents = nil
		newEvents = append(newEvents, &errorEvent)
		err = nil
	}
	events := make([]EventRecord, len(newEvents))
	for i := 0; i < len(newEvents); i++ {
		newEvents[i].SetID(lib.MustNewULID())
		newEvents[i].SetAggregateID(rec.AggregateID)
		newEvents[i].SetEventType(reflect.TypeOf(newEvents[i]).Elem().Name())
		newEvents[i].SetVersion(expectedVersion + i)
		events[i], err = EventToEventRecord(newEvents[i])
		if err != nil {
			return err
		}
	}
	err = c.store.StoreCommandResults(ctx, rec.ID, expectedVersion, events...)

	return
}
