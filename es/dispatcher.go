package es

import "context"

type CommandDispatcher interface {
	DispatchCommandRequest(ctx context.Context, request CommandRequest) (string, error)
	DispatchCommand(ctx context.Context, command ICommand) (string, error)
	Close()
}
