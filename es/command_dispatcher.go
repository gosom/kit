package es

import "context"

type CommandDispatcher interface {
	Dispatch(ctx context.Context, cmd ICommand) (string, error)
}
