package kafka

import (
	"context"
	"fmt"
	"sync"
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
	commitWg    *sync.WaitGroup
}

func NewConsumer(topic string, commitEvery int, cfg kafka.ConfigMap, w es.Worker) (*Consumer, error) {
	ans := Consumer{
		log:         logging.Get().With("component", "kafka", "topic", topic),
		topic:       topic,
		offsetsMap:  make(map[string]kafka.TopicPartition),
		worker:      w,
		commitEvery: commitEvery,
		commitWg:    &sync.WaitGroup{},
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
	defer func() {
		commit(o.consumer, o.offsetsMap, o.log.Error, o.commitWg)
		o.commitWg.Wait()
		time.Sleep(2 * time.Second)
		o.consumer.Close()
		o.log.Info("Consumer closed")
	}()
	for {
		select {
		case <-ctx.Done():
			go commit(o.consumer, o.offsetsMap, o.log.Error, o.commitWg)
			return nil
		default:
		}
		msg, err := o.consumer.ReadMessage(100 * time.Millisecond)
		if err == nil {
			o.log.Debug("Received message", "key", string(msg.Key), "value", string(msg.Value))

			// we process the message here
			if err := o.processMessage(ctx, msg); err != nil {
				if err == context.Canceled {
					o.log.Info("Context canceled")
					return nil
				}
				panic(err)
			}
			fmt.Println(msg.TopicPartition.Offset)

			key := fmt.Sprintf("%s[%d]", *msg.TopicPartition.Topic, msg.TopicPartition.Partition)
			o.offsetsMap[key] = msg.TopicPartition

			o.count++
			if o.count%o.commitEvery == 0 {
				go commit(o.consumer, o.offsetsMap, o.log.Error, o.commitWg)
			}
		}
	}
	return nil
}

// processMessage is the place where we process the message
func (o *Consumer) processMessage(ctx context.Context, msg *kafka.Message) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %v", r)
		}
	}()
	backoff := 20 * time.Millisecond
	factor := 2
	maxWait := 5 * time.Second
	for {
		err = o.worker.Process(ctx, msg.Key, msg.Value, msg.Timestamp)
		//err = fmt.Errorf("artificial error")
		if err == nil {
			return
		}
		o.log.Error("Error processing message", "error", err, "func", "processMessage")
		select {
		case <-ctx.Done():
			o.log.Info("Context canceled", "func", "processMessage")
			return ctx.Err()
		default:
			if backoff > maxWait {
				backoff = maxWait
			}
			o.log.Info("Retrying in", "backoff", backoff, "func", "processMessage")
			time.Sleep(backoff)
			backoff = backoff * time.Duration(factor)
		}
	}
	return
}

func (o *Consumer) rebalanceCb(consumer *kafka.Consumer, ev kafka.Event) error {
	switch e := ev.(type) {
	case kafka.AssignedPartitions:
		o.count = 0
		o.offsetsMap = make(map[string]kafka.TopicPartition)
		consumer.Assign(e.Partitions)
		o.log.Info("RebalanceCb - AssignedPartitions", "partitions", e.Partitions)
	case kafka.RevokedPartitions:
		commit(consumer, o.offsetsMap, o.log.Error, o.commitWg)
		consumer.Unassign()
		o.log.Info("RebalanceCb - RevokedPartitions", "partitions", e.Partitions)
	}
	return nil
}

func commit(consumer *kafka.Consumer, offsets map[string]kafka.TopicPartition, logFn func(string, ...any), wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()
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
	logFn("Offsets committed", "offsets", tps)
}
