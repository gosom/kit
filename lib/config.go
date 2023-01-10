package lib

import "github.com/kelseyhightower/envconfig"

func NewConfig[C any](prefix string) (C, error) {
	var cfg C
	err := envconfig.Process(prefix, &cfg)
	return cfg, err
}
