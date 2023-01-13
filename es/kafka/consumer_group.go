package kafka

import (
	"context"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gosom/kit/es"
	"golang.org/x/sync/errgroup"
)

type ConsumerGroup struct {
	cfg         kafka.ConfigMap
	topic       string
	num         int
	worker      es.Worker
	commitEvery int
}

func NewConsumerGroup(cfg KafkaConfig, topic string, num int, w es.Worker) *ConsumerGroup {
	ans := ConsumerGroup{
		cfg:         NewKafkaConfigMap(cfg),
		topic:       topic,
		num:         num,
		worker:      w,
		commitEvery: 100,
	}
	return &ans
}

func (o *ConsumerGroup) Listen(ctx context.Context) error {
	consumers := make([]*Consumer, 0, o.num)
	for i := 0; i < o.num; i++ {
		c, err := NewConsumer(o.topic, o.commitEvery, o.cfg, o.worker)
		if err != nil {
			return err
		}
		consumers = append(consumers, c)
	}
	g, ctx := errgroup.WithContext(ctx)
	for i := range consumers {
		consumer := consumers[i]
		g.Go(func() error {
			return consumer.Start(ctx)
		})
	}
	return g.Wait()
}
