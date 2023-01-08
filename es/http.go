package es

import (
	"context"
)

type WebServer interface {
	ListenAndServe(ctx context.Context) error
}
