package lib

import "context"

type contextKey int

const (
	errorKey     contextKey = 1
	clientIPKey  contextKey = 2
	userKey      contextKey = 3
	requestIDKey contextKey = 4
)

func NewContextWithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, requestIDKey, id)
}

func RequestIDFromContext(ctx context.Context) string {
	id, ok := ctx.Value(requestIDKey).(string)
	if !ok {
		return ""
	}
	return id
}

func NewContextWithUser(ctx context.Context, user User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

func UserFromContext(ctx context.Context) User {
	user, ok := ctx.Value(userKey).(User)
	if !ok {
		return nil
	}
	return user
}

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
