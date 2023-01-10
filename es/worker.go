package es

import (
	"context"
	"time"
)

type Worker interface {
	Process(ctx context.Context, key, value []byte, timestamp time.Time) error
}
