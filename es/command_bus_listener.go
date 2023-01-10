package es

import "context"

type CommandBusListener interface {
	Listen(ctx context.Context) error
}
