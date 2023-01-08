package kafka

/*  This is WIP

import (
	"context"
	"fmt"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gosom/kit/es"
)

type Consumer struct {
	topic      string
	consumer   *kafka.Consumer
	offsetsMap map[string]kafka.TopicPartition
	count      int
	worker     es.Worker
}

func NewConsumer(topic string, cfg kafka.ConfigMap, w es.Worker) (*Consumer, error) {
	ans := Consumer{
		topic:      topic,
		offsetsMap: make(map[string]kafka.TopicPartition),
		worker:     w,
	}
	consumer, err := kafka.NewConsumer(&cfg)
	if err != nil {
		return nil, err
	}
	ans.consumer = consumer
	return &ans, nil
}

func (o *Consumer) Start(ctx context.Context) error {
	o.consumer.SubscribeTopics([]string{o.topic}, o.rebalanceCb)
	for {
		select {
		case <-ctx.Done():
			commit(o.consumer, o.offsetsMap)
			o.consumer.Close()
			time.Sleep(2 * time.Second)
			fmt.Println("Consumer - Done")
			return nil
		default:
			msg, err := o.consumer.ReadMessage(100 * time.Millisecond)
			if err == nil {
				fmt.Printf("Consumer - Message on %s: %s\n", msg.TopicPartition, string(msg.Value))

				if err := o.worker.Process(ctx, msg.Key, msg.Value, msg.Timestamp); err != nil {
					return err
				}

				key := fmt.Sprintf("%s[%d]", *msg.TopicPartition.Topic, msg.TopicPartition.Partition)
				o.offsetsMap[key] = msg.TopicPartition

				o.count++
				if o.count%10 == 0 {
					go commit(o.consumer, o.offsetsMap)
				}
			}
		}
	}
	return nil
}

func (o *Consumer) rebalanceCb(consumer *kafka.Consumer, ev kafka.Event) error {
	switch e := ev.(type) {
	case kafka.AssignedPartitions:
		fmt.Println("Rebalance - Assigned:", e.Partitions)

		// Reset the state
		o.count = 0
		o.offsetsMap = make(map[string]kafka.TopicPartition)

		// Assign partition
		consumer.Assign(e.Partitions)
	case kafka.RevokedPartitions:
		fmt.Println("Rebalance - Revoked:", e.Partitions)
		// Commit the current offset synchronously before revoked partitions
		commit(consumer, o.offsetsMap)

		consumer.Unassign()
	}
	return nil
}

func commit(consumer *kafka.Consumer, offsets map[string]kafka.TopicPartition) {
	if len(offsets) == 0 {
		return
	}
	tps := make([]kafka.TopicPartition, len(offsets))
	index := 0
	for _, tp := range offsets {
		// The committed offset should always be the offset of the next message that your application will read.
		tp.Offset = tp.Offset + 1
		tps[index] = tp
		index++
	}
	if _, err := consumer.CommitOffsets(tps); err != nil {
		fmt.Println("CommitOffsets Error:", err)
		return
	}
}
*/
