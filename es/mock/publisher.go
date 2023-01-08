package mock

import (
	"context"
	"fmt"

	"github.com/gosom/kit/es"
)

var _ es.Publisher = (*Publisher)(nil)

type Publisher struct {
}

func NewPublisher() (*Publisher, error) {
	return &Publisher{}, nil
}

func (p *Publisher) Name() string {
	return "mock"
}

func (p *Publisher) Publish(ctx context.Context, events ...es.EventRecord) error {
	for i := range events {
		fmt.Println(events[i])
	}
	return nil
}
