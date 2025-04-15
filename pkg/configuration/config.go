package configuration

import (
	"errors"

	"github.com/caarlos0/env/v11"
)

var (
	ErrFailedToParseConfig = errors.New("failed to parse config from env")
)

func Parse[T any]() (*T, error) {
	var cfg T
	if err := env.Parse(&cfg); err != nil {
		return &cfg, ErrFailedToParseConfig
	}
	return &cfg, nil
}
