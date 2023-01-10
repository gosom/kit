package es

import (
	"context"
	"errors"

	"github.com/gosom/kit/logging"
	"golang.org/x/sync/errgroup"
)

type option func(*appService) error

type appService struct {
	log logging.Logger

	store              EventStore
	webServer          WebServer
	commandProcessor   CommandProcessor
	commandBusListener CommandBusListener

	subscribers []Subscriber
}

func (a *appService) Start(ctx context.Context) error {
	defer func() {
		a.log.Info("application service stopped")
	}()
	a.log.Info("starting application service")
	g, ctx := errgroup.WithContext(ctx)
	for i := range a.subscribers {
		sub := a.subscribers[i]
		g.Go(func() error {
			return sub.Start(ctx)
		})
	}
	if a.commandProcessor != nil {
		g.Go(func() error {
			return a.commandProcessor.Start(ctx)
		})
	}
	if a.webServer != nil {
		g.Go(func() error {
			return a.webServer.ListenAndServe(ctx)
		})
	}
	if a.commandBusListener != nil {
		g.Go(func() error {
			return a.commandBusListener.Listen(ctx)
		})
	}
	return g.Wait()
}

func New(options ...option) (*appService, error) {
	app := appService{}
	for _, opt := range options {
		if err := opt(&app); err != nil {
			return nil, err
		}
	}
	return &app, nil
}

func WithLogger(log logging.Logger) option {
	return func(a *appService) error {
		a.log = log
		return nil
	}
}

func WithEventStore(store EventStore) option {
	return func(a *appService) error {
		a.store = store
		return nil
	}
}

func WithWebServer(srv WebServer) option {
	return func(a *appService) error {
		a.webServer = srv
		if a.store == nil {
			return errors.New("event store is not set")
		}
		return nil
	}
}

func WithCommandProcessor(processor CommandProcessor) option {
	return func(a *appService) error {
		a.commandProcessor = processor
		return nil
	}
}

func WithPublishers(publisher ...Publisher) option {
	return func(a *appService) error {
		if a.store == nil {
			return errors.New("event store is not set")
		}
		for i := range publisher {
			sub, err := NewSubscriber(a.store, publisher[i], publisher[i].Name())
			if err != nil {
				return err
			}
			a.subscribers = append(a.subscribers, sub)
		}
		return nil
	}
}

func WithCommandBusListener(listener CommandBusListener) option {
	return func(a *appService) error {
		a.commandBusListener = listener
		return nil
	}
}
