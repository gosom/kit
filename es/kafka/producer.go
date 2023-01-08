package kafka

/* WIP
import (
	"encoding/json"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gosom/kit/es"
	"github.com/gosom/kit/logging"
)

type KafkaProducer struct {
	p *kafka.Producer
	l logging.Logger
}

func NewKafkaProducer(servers string) (*KafkaProducer, error) {
	l := logging.Get().With("component", "kafka")
	cfg := make(kafka.ConfigMap)
	cfg.SetKey("bootstrap.servers", "")
	cfg.SetKey("security.protocol", "SASL_SSL")
	cfg.SetKey("sasl.mechanisms", "PLAIN")
	cfg.SetKey("sasl.username", "")
	cfg.SetKey("sasl.password", "")
	//if err := cfg.SetKey("bootstrap.servers", servers); err != nil {
	//	return nil, err
	//}
	if err := cfg.SetKey("go.delivery.reports", true); err != nil {
		return nil, err
	}
	cfg.SetKey("go.logs.channel.enable", true)
	cfg.SetKey("api.version.request", false)
	cfg.SetKey("compression.type", "snappy")
	p, err := kafka.NewProducer(&cfg)
	if err != nil {
		return nil, err
	}
	return &KafkaProducer{p: p, l: l}, nil
}

func (o *KafkaProducer) Dispatch(item es.MessageEncoder) error {
	msg, err := item.Encode()
	if err != nil {
		return err
	}
	topic := msg.Topic
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	ch := make(chan kafka.Event, 1)
	defer close(ch)
	err = o.p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            []byte(msg.AggregateID),
		Value:          data,
	}, ch)
	if err != nil {
		return err
	}
	e := <-ch
	switch ev := e.(type) {
	case *kafka.Message:
		if ev.TopicPartition.Error != nil {
			return ev.TopicPartition.Error
		}
		o.l.Info("dispatched", "message", msg.Type, "topic", topic, "key", msg.AggregateID, "offset", ev.TopicPartition.Offset)
		return nil
	default:
		return fmt.Errorf("unknown event type: %T", ev)
	}
	return nil
}

func (o *KafkaProducer) Close() {
	o.p.Close()
}
*/
