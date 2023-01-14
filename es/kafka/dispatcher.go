package kafka

import (
	"bytes"
	"context"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gosom/kit/es"
	"github.com/gosom/kit/lib"
	"github.com/gosom/kit/logging"
)

type Dispatcher struct {
	log      logging.Logger
	domain   string
	ack      bool
	topic    string
	p        *kafka.Producer
	registry *es.Registry
}

func NewDispatcher(cfg kafka.ConfigMap, ack bool, topic, domain string, registry *es.Registry) (*Dispatcher, error) {
	ans := Dispatcher{
		log:      logging.Get().With("component", "kafka_dispatcher"),
		domain:   domain,
		topic:    topic,
		registry: registry,
		ack:      ack,
	}
	var err error
	ans.p, err = kafka.NewProducer(&cfg)
	if err != nil {
		return nil, err
	}
	return &ans, nil
}

func (d *Dispatcher) DispatchCommandRequest(ctx context.Context, request es.CommandRequest) (string, error) {
	if err := lib.Validate(request); err != nil {
		return "", fmt.Errorf("%w %s", es.ErrInvalidCommand, err.Error())
	}
	r := bytes.NewReader(request.Payload)
	command, err := es.ParseCommandRequest(d.registry, r)
	if err != nil {
		return "", nil
	}
	return d.DispatchCommand(ctx, command)
}

func (d *Dispatcher) DispatchCommand(ctx context.Context, command es.ICommand) (string, error) {
	cr, err := es.CommandToCommandRecord(d.domain, command)
	if err != nil {
		return "", err
	}
	msg, err := es.CommandRecordToBusMessage(cr)
	if err != nil {
		return "", err
	}
	busmsg := kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &d.topic,
			Partition: kafka.PartitionAny,
		},
		Key:   []byte(cr.AggregateID),
		Value: msg.Data,
	}
	switch d.ack {
	case true:
		err = d.produceWithAck(ctx, &busmsg)
	default:
		err = d.produceWithoutAck(ctx, &busmsg)
	}
	if err != nil {
		return "", err
	}
	return cr.ID, nil
}

func (d *Dispatcher) produceWithAck(ctx context.Context, msg *kafka.Message) error {
	ch := make(chan kafka.Event, 1)
	defer close(ch)
	if err := d.p.Produce(msg, ch); err != nil {
		return err
	}
	e := <-ch
	switch ev := e.(type) {
	case *kafka.Message:
		if ev.TopicPartition.Error != nil {
			return ev.TopicPartition.Error
		}
		return nil
	default:
		return fmt.Errorf("unknown event type: %T", ev)
	}
	return nil
}
func (d *Dispatcher) produceWithoutAck(ctx context.Context, msg *kafka.Message) error {
	if err := d.p.Produce(msg, nil); err != nil {
		return err
	}
	return nil
}

func (d *Dispatcher) Close() {
	d.p.Flush(2000)
	d.p.Close()
}
