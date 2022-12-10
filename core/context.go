package core

import "context"

type contextKey int

const (
	errorKey    contextKey = 1
	clientIPKey contextKey = 2
)

func NewContextWithErr(ctx context.Context, err error) context.Context {
	return context.WithValue(ctx, errorKey, err)
}

func ErrorFromContext(ctx context.Context) error {
	err, ok := ctx.Value(errorKey).(error)
	if !ok {
		return nil
	}
	return err
}

func NewContextWithClientIP(ctx context.Context, ip string) context.Context {
	return context.WithValue(ctx, clientIPKey, ip)
}

func IPFromContext(ctx context.Context) string {
	ip, ok := ctx.Value(clientIPKey).(string)
	if !ok {
		return ""
	}
	return ip
}
