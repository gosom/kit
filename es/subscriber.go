package es

import (
	"context"
	"fmt"
	"time"

	"github.com/gosom/kit/logging"
)

type Subscriber interface {
	Start(ctx context.Context) error
}

var _ Subscriber = (*subscriber)(nil)

type subscriber struct {
	publisher    Publisher
	store        EventStore
	subscription Subscription
	log          logging.Logger
}

func NewSubscriber(store EventStore, publisher Publisher, subscription string) (Subscriber, error) {
	sub, err := store.InsertSubscription(context.Background(), subscription)
	if err != nil {
		return nil, err
	}
	ans := subscriber{
		publisher:    publisher,
		store:        store,
		subscription: sub,
		log:          logging.Get().With("component", "es/subscriber"),
	}
	return &ans, nil
}

func (o *subscriber) Start(ctx context.Context) error {
	o.log.Info("starting subscriber", "subscription", o.subscription.Group)
	defer o.log.Info("subscriber stopped", "subscription", o.subscription.Group)
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			if num, err := o.process(ctx); err != nil {
				o.log.Error("Error processing events", "subscription", o.subscription.Group, "error", err)
			} else if num > 0 {
				o.log.Info("Processed events", "subscription", o.subscription.Group, "num", num)
			}
		}
	}
}

func (o *subscriber) process(ctx context.Context) (int, error) {
	items, err := o.store.SelectEventsForSubscription(ctx, o.subscription, 100)
	if err != nil {
		return 0, fmt.Errorf("%w when selecting events for subscription", err)
	}
	if len(items) > 0 {
		if err := o.publisher.Publish(ctx, items...); err != nil {
			return 0, fmt.Errorf("%w when publishing events", err)
		}
		o.subscription, err = o.store.UpdateSubscription(ctx, o.subscription.Group, items[len(items)-1].ID)
		if err != nil {
			return 0, fmt.Errorf("%w when updating subscription", err)
		}
	}
	return len(items), nil
}
