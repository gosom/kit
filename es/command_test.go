package es_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/gosom/kit/es"
	"github.com/stretchr/testify/require"
)

type problematicCommand struct {
	es.CommandBase
	ID bool `json:"id" validate:"required" aggregateID:"true"`
}

type dummyCommand struct {
	es.CommandBase
	ID    string `json:"id" validate:"required,uuid" aggregateID:"true"`
	Title string `json:"title" validate:"required"`
}

type dummyLoader struct {
}

func (l *dummyLoader) Load(ctx context.Context, aggregateID string, agg es.AggregateRoot) error {
	return nil
}

func TestCommandBase(t *testing.T) {
	t.Run("NewCommandBase", func(t *testing.T) {
		cb := es.CommandBase{}
		cb.SetID("123")
		require.Equal(t, "123", cb.GetID())
		cb.SetEventType("test")
		require.Equal(t, "test", cb.GetEventType())
		cb.SetAggregateID("test-123")
		require.Equal(t, "test-123", cb.GetAggregateID())
		cb.SetAggregateHash()
		require.Equal(t, int32(1442536664), cb.GetAggregateHash())
		err := cb.Validate()
		require.NoError(t, err)
	})
	t.Run("CommandBaseValidateWithErr", func(t *testing.T) {
		cb := es.CommandBase{}

		err := cb.Validate()
		require.Error(t, err)
		require.ErrorIs(t, err, es.ErrInvalidCommand)

		cb.SetID("123")
		err = cb.Validate()
		require.Error(t, err)
		require.ErrorIs(t, err, es.ErrInvalidCommand)

		cb.SetEventType("test")
		err = cb.Validate()
		require.Error(t, err)
		require.ErrorIs(t, err, es.ErrInvalidCommand)

		cb.SetAggregateID("123")
		err = cb.Validate()
		require.Error(t, err)
		require.ErrorIs(t, err, es.ErrInvalidCommand)

		cb.SetAggregateID("test-123")
		err = cb.Validate()
		require.Error(t, err)
		require.ErrorIs(t, err, es.ErrInvalidCommand)

		cb.SetAggregateHash()
		err = cb.Validate()
		require.NoError(t, err)
	})
	t.Run("CommandBase test Handle", func(t *testing.T) {
		cb := es.CommandBase{}
		cb.SetID("123")
		require.Equal(t, "123", cb.GetID())
		cb.SetEventType("test")
		require.Equal(t, "test", cb.GetEventType())
		cb.SetAggregateID("test-123")
		require.Equal(t, "test-123", cb.GetAggregateID())
		cb.SetAggregateHash()
		require.Equal(t, int32(1442536664), cb.GetAggregateHash())
		err := cb.Validate()
		require.NoError(t, err)

		loader := &dummyLoader{}
		require.Implements(t, (*es.AggregateLoader)(nil), loader)

		err = func() (err error) {
			defer func() {
				if r := recover(); r != nil {
					err = fmt.Errorf("panic: %v", r)
				}
			}()
			_, err = cb.Handle(context.Background(), loader)
			return err
		}()
		require.Error(t, err)
		require.ErrorContains(t, err, "panic: not implemented")
	})
}

func TestCommandToCommandRecord(t *testing.T) {
	t.Run("Test with CommandBase", func(t *testing.T) {
		cb := es.CommandBase{}
		cb.SetID("123")
		cb.SetEventType("test")
		require.Equal(t, "test", cb.GetEventType())
		cb.SetAggregateID("test-123")
		require.Equal(t, "test-123", cb.GetAggregateID())
		cb.SetAggregateHash()
		require.Equal(t, int32(1442536664), cb.GetAggregateHash())
		err := cb.Validate()
		require.NoError(t, err)

		cr, err := es.CommandToCommandRecord("test", &cb)
		require.NoError(t, err)
		require.Equal(t, "123", cr.ID)
		require.Equal(t, "test", cr.EventType)
		require.Equal(t, "test-123", cr.AggregateID)
		require.Equal(t, int32(1442536664), cr.AggregateHash)

		require.JSONEq(t, `{}`, string(cr.Data))
	})

	t.Run("Test with Dummy Command", func(t *testing.T) {
		cb := dummyCommand{}

		cr, err := es.CommandToCommandRecord("test", &cb)
		require.Error(t, err)
		require.ErrorIs(t, err, es.ErrInvalidCommand)

		cb.ID = "123" // this is not valid we need uuid
		cb.Title = "test"
		cr, err = es.CommandToCommandRecord("test", &cb)
		require.Error(t, err)
		require.ErrorIs(t, err, es.ErrInvalidCommand)

		cb.ID = "8f6ee4b0-9970-11ed-918c-4b15b2f0ca00"
		cr, err = es.CommandToCommandRecord("test", &cb)
		require.NoError(t, err)
		require.IsType(t, es.CommandRecord{}, cr)

		require.NotEqual(t, "8f6ee4b0-9970-11ed-918c-4b15b2f0ca00", cr.ID)
		require.Len(t, cr.ID, 26)
		require.Equal(t, "dummyCommand", cr.EventType)
		require.Greater(t, cr.AggregateHash, int32(0))
		require.Greater(t, cb.GetAggregateHash(), int32(0))
	})
	t.Run("Test with problematic Command", func(t *testing.T) {
		cb := problematicCommand{}

		cb.ID = true
		_, err := es.CommandToCommandRecord("test", &cb)
		require.Error(t, err)
		require.ErrorIs(t, err, es.ErrInvalidEvent)
		require.ErrorContains(t, err, "aggregateID field must be of type string, int or uint")

	})
}
