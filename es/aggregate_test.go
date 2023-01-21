package es_test

import (
	"fmt"
	"testing"

	"github.com/gosom/kit/es"
	"github.com/stretchr/testify/require"
)

type dummyEventWithErr struct {
	es.EventBase
}

func (e *dummyEventWithErr) Apply(agg es.AggregateRoot) error {
	return fmt.Errorf("dummy error")
}

type dummyEvent struct {
	es.EventBase
}

func (e *dummyEvent) Apply(agg es.AggregateRoot) error {
	return nil
}

func TestAggregate(t *testing.T) {
	t.Run("NewAggregate", func(t *testing.T) {
		var factory es.AggregateFactory = func() (es.AggregateRoot, error) {
			return es.NewAggregateBase()
		}
		require.NotNil(t, factory)
		agg, err := factory()
		require.NoError(t, err)
		require.NotNil(t, agg)
		require.IsType(t, &es.AggregateBase{}, agg)
		require.Implements(t, (*es.AggregateRoot)(nil), agg)

		require.Equal(t, "", agg.GetID())
		agg.SetID("123")
		require.Equal(t, "123", agg.GetID())

		require.Equal(t, uint64(0), agg.GetVersion())
		agg.SetVersion(1)
		require.Equal(t, uint64(1), agg.GetVersion())

		require.Equal(t, "", agg.GetType())
		agg.SetType("test")
		require.Equal(t, "test", agg.GetType())

		require.Equal(t,
			"AggregateBase{ID: 123, Type: test, Version: 1}",
			agg.String(),
		)
	})
	t.Run("Test Raise/Load Events", func(t *testing.T) {
		var factory es.AggregateFactory = func() (es.AggregateRoot, error) {
			return es.NewAggregateBase()
		}
		require.NotNil(t, factory)
		agg, err := factory()
		require.NoError(t, err)
		require.NotNil(t, agg)
		require.IsType(t, &es.AggregateBase{}, agg)
		require.Implements(t, (*es.AggregateRoot)(nil), agg)

		require.Equal(t, uint64(0), agg.GetVersion())

		ev := &dummyEvent{}
		ev.SetVersion(1)

		es.RaiseEvent(agg, ev)
		require.Equal(t, uint64(1), agg.GetVersion())

		ev2 := &dummyEvent{}
		ev2.SetVersion(2)

		ev3 := &dummyEvent{}
		ev3.SetVersion(3)

		es.Load(agg, []es.IEvent{ev2, ev3})
		require.Equal(t, uint64(3), agg.GetVersion())

		// now let's add a bad event
		err = es.RaiseEvent(agg, nil)
		require.Error(t, err)
		require.ErrorIs(t, err, es.ErrInvalidEvent)

		// now lets add an event seen previously
		err = es.RaiseEvent(agg, ev3)
		require.Error(t, err)
		require.ErrorIs(t, err, es.ErrDuplicateEvent)

		// let's try to Raise event with nil aggregate
		err = es.RaiseEvent(nil, ev3)
		require.Error(t, err)
		require.ErrorIs(t, err, es.ErrNilAggregate)

		// let's try to Raise event that Apply returns error
		ev5 := &dummyEventWithErr{}
		ev5.SetVersion(5)
		err = es.RaiseEvent(agg, ev5)
		require.Error(t, err)
		require.Equal(t, "dummy error", err.Error())

		// let's try to Load events with nil aggregate
		err = es.Load(nil, []es.IEvent{ev2, ev3})
		require.Error(t, err)
	})
}
