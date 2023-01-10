package kafka

import "context"

type CommandBusListener struct {
}

func (c *CommandBusListener) Listen(ctx context.Context) error {
	panic("not implemented") // TODO: Implement
}
