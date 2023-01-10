package kafka

import (
	"context"
	"fmt"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gosom/kit/es"
	"github.com/gosom/kit/logging"
)

type Consumer struct {
	log         logging.Logger
	topic       string
	consumer    *kafka.Consumer
	offsetsMap  map[string]kafka.TopicPartition
	count       int
	worker      es.Worker
	commitEvery int
}

func NewConsumer(topic string, commitEvery int, cfg kafka.ConfigMap, w es.Worker) (*Consumer, error) {
	ans := Consumer{
		log:         logging.Get().With("component", "kafka", "topic", topic),
		topic:       topic,
		offsetsMap:  make(map[string]kafka.TopicPartition),
		worker:      w,
		commitEvery: commitEvery,
	}
	consumer, err := kafka.NewConsumer(&cfg)
	if err != nil {
		return nil, err
	}
	ans.consumer = consumer
	return &ans, nil
}

func (o *Consumer) Start(ctx context.Context) error {
	o.log.Info("Starting consumer")
	o.consumer.SubscribeTopics([]string{o.topic}, o.rebalanceCb)
	for {
		select {
		case <-ctx.Done():
			commit(o.consumer, o.offsetsMap, o.log.Error)
			o.consumer.Close()
			time.Sleep(2 * time.Second)
			o.log.Info("Consumer closed")
			return nil
		default:
		}
		msg, err := o.consumer.ReadMessage(100 * time.Millisecond)
		if err == nil {
			o.log.Debug("Received message", "key", string(msg.Key), "value", string(msg.Value))

			// we process the message here
			if err := o.processMessage(ctx, msg); err != nil {
				if err == context.Canceled {
					continue
				}
				panic(err)
			}

			key := fmt.Sprintf("%s[%d]", *msg.TopicPartition.Topic, msg.TopicPartition.Partition)
			o.offsetsMap[key] = msg.TopicPartition

			o.count++
			if o.count%o.commitEvery == 0 {
				o.log.Info("Committing offsets", "offsets", o.offsetsMap)
				go commit(o.consumer, o.offsetsMap, o.log.Error)
			}
		}
	}
	return nil
}

// processMessage is the place where we process the message
func (o *Consumer) processMessage(ctx context.Context, msg *kafka.Message) error {
	backoff := 20 * time.Millisecond
	factor := 2
	maxWait := 5 * time.Second
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		err := o.worker.Process(ctx, msg.Key, msg.Value, msg.Timestamp)
		if err == nil {
			return nil
		}
		o.log.Error("Error processing message", "error", err)
		if backoff > maxWait {
			backoff = maxWait
		}
		time.Sleep(backoff)
		backoff = backoff * time.Duration(factor)
	}
	return nil
}

func (o *Consumer) rebalanceCb(consumer *kafka.Consumer, ev kafka.Event) error {
	switch e := ev.(type) {
	case kafka.AssignedPartitions:
		o.log.Info("RebalanceCb - AssignedPartitions", "partitions", e.Partitions)

		o.count = 0
		o.offsetsMap = make(map[string]kafka.TopicPartition)

		consumer.Assign(e.Partitions)
	case kafka.RevokedPartitions:
		o.log.Info("RebalanceCb - RevokedPartitions", "partitions", e.Partitions)

		commit(consumer, o.offsetsMap, o.log.Error)

		consumer.Unassign()
	}
	return nil
}

func commit(consumer *kafka.Consumer, offsets map[string]kafka.TopicPartition, logFn func(string, ...any)) {
	if len(offsets) == 0 {
		return
	}
	tps := make([]kafka.TopicPartition, len(offsets))
	index := 0
	for _, tp := range offsets {
		tp.Offset = tp.Offset + 1
		tps[index] = tp
		index++
	}
	if _, err := consumer.CommitOffsets(tps); err != nil {
		logFn("Error committing offsets", "error", err)
		return
	}
}
