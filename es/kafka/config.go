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
	Producer          bool
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

	m := kafka.ConfigMap{
		"bootstrap.servers":   cfg.Servers,
		"security.protocol":   cfg.Security,
		"sasl.mechanisms":     cfg.Mechanism,
		"sasl.username":       cfg.Username,
		"sasl.password":       cfg.Password,
		"api.version.request": cfg.ApiVersionRequest,
	}
	if cfg.Producer {
		m.SetKey("enable.idempotence", true)
	} else {
		m.SetKey("enable.auto.commit", false)
		m.SetKey("go.application.rebalance.enable", true)
		m.SetKey("enable.auto.commit", true)
		m.SetKey("group.id", cfg.GroupID)
		m.SetKey("auto.offset.reset", cfg.AutoOffsetReset)
	}
	return m
}

func NewKafkaConfig() (KafkaConfig, error) {
	cfg, err := lib.NewConfig[KafkaConfig]("")
	if err != nil {
		return KafkaConfig{}, err
	}
	return cfg, nil
}
