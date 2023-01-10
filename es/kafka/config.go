package kafka

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gosom/kit/lib"
)

type KafkaConfig struct {
	Servers           string
	Security          string
	Mechanism         string
	Username          string
	Password          string
	GroupID           string
	RebalanceEnable   bool
	AutoOffsetReset   string
	ApiVersionRequest bool
}

func NewKafkaConfigMap(cfg KafkaConfig) kafka.ConfigMap {
	if cfg.Security == "" {
		cfg.Security = "PLAINTEXT"
	}
	if cfg.Mechanism == "" {
		cfg.Mechanism = "PLAIN"
	}
	if cfg.GroupID == "" {
		cfg.GroupID = "default"
	}
	if cfg.AutoOffsetReset == "" {
		cfg.AutoOffsetReset = "latest"
	}

	return kafka.ConfigMap{
		"bootstrap.servers":               cfg.Servers,
		"security.protocol":               cfg.Security,
		"sasl.mechanisms":                 cfg.Mechanism,
		"sasl.username":                   cfg.Username,
		"sasl.password":                   cfg.Password,
		"group.id":                        cfg.GroupID,
		"enable.auto.commit":              false,
		"go.application.rebalance.enable": true,
		"auto.offset.reset":               cfg.AutoOffsetReset,
		"api.version.request":             cfg.ApiVersionRequest,
	}
}

func NewKafkaConfig() (KafkaConfig, error) {
	cfg, err := lib.NewConfig[KafkaConfig]("")
	if err != nil {
		return KafkaConfig{}, err
	}
	return cfg, nil
}
