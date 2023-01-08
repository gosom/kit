package kafka

/* This is WIP

import (
	"context"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gosom/kit/es"
	"golang.org/x/sync/errgroup"
)

const (
	servers = ""
	securit = "SASL_SSL"
	mechani = "PLAIN"
	usernam = ""
	passwor = ""
)

type ConsumerGroup struct {
	Name    string
	Brokers []string
	GroupID string
	Num     int
	Worker  es.Worker
	Topic   string
}

func (o *ConsumerGroup) Start(ctx context.Context) error {
	cfg := kafka.ConfigMap{
		"bootstrap.servers":               servers,
		"security.protocol":               securit,
		"sasl.mechanisms":                 mechani,
		"sasl.username":                   usernam,
		"sasl.password":                   passwor,
		"group.id":                        o.GroupID,
		"enable.auto.commit":              false,
		"go.application.rebalance.enable": true,
		"auto.offset.reset":               "latest",
		"api.version.request":             false,
	}

	consumers := make([]*Consumer, 0, o.Num)
	for i := 0; i < o.Num; i++ {
		c, err := NewConsumer(topic, cfg, o.Worker)
		if err != nil {
			return err
		}
		consumers = append(consumers, c)
	}
	g, ctx := errgroup.WithContext(ctx)
	for i := range consumers {
		g.Go(func() error {
			return consumers[i].Start(ctx)
		})
	}
	return g.Wait()
}
*/
