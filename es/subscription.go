package es

import (
	"time"
)

type Subscription struct {
	Group           string
	LastSeenEventID string
	LastUpdatedAt   time.Time
}

func (o *Subscription) Bind() []any {
	return []any{&o.Group, &o.LastSeenEventID, &o.LastUpdatedAt}
}
