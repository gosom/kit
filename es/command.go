package es

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strings"
	"time"

	"github.com/gosom/kit/core"
	"github.com/gosom/kit/lib"
)

var _ ICommand = (*CommandBase)(nil)

// CommandBase is the base struct for commands.
type CommandBase struct {
	CommandEventBase
	aggregateHash int32
}

func (c *CommandBase) GetAggregateHash() int32 {
	return c.aggregateHash
}

func (c *CommandBase) SetAggregateHash() {
	c.aggregateHash = lib.Int32Ring(lib.HashToUInt32(c.aggregateID))
}

func (c *CommandBase) Handle(ctx context.Context, h AggregateLoader) ([]IEvent, error) {
	panic("not implemented")
}

// validate validates the event.
func (e *CommandBase) Validate() error {
	if len(e.id) == 0 {
		return fmt.Errorf("Command ID is required %w", ErrInvalidCommand)
	}
	if len(e.aggregateID) == 0 {
		return fmt.Errorf("Command AggregateID is required %w", ErrInvalidCommand)
	}
	before, after, ok := strings.Cut(e.aggregateID, "-")
	if !ok || len(before) == 0 || len(after) == 0 {
		return fmt.Errorf("Command AggregateID is invalid %w", ErrInvalidCommand)
	}

	if e.aggregateHash == 0 {
		return fmt.Errorf("Command AggregateHash is required %w", ErrInvalidCommand)
	}
	if len(e.eventType) == 0 {
		return fmt.Errorf("Command Event Type is required %w", ErrInvalidCommand)
	}
	return nil
}

// CommandRecord is the record for a command.
type CommandRecord struct {
	RecordBase
	AggregateHash int32
	Status        string
}

func (o *CommandRecord) Bind() []any {
	ans := o.RecordBase.Bind()
	ans = append(ans, &o.AggregateHash, &o.Status)
	return ans
}

// prepareCommand sets the command ID and the event type and the aggregate ID.
// It also validates the command.
// this method is called before a command is published to the command bus.
// When the id is not set, a new ULID is generated.
// When the event type is not set, the event type is set to the command's struct name.
// When the aggregate ID is not set, the aggregate ID is set to the value of the field with the tag `aggregateID:"true"`.
func prepareCommand(domain string, ev ICommand) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("error preparing command %+v: %v", ev, r)
		}
	}()

	if err := core.Validate(ev); err != nil {
		return fmt.Errorf("%s %w", err, ErrInvalidCommand)
	}

	if ev.GetID() == "" {
		ev.SetID(lib.MustNewULID())
	}
	v := reflect.ValueOf(ev)
	elem := v.Elem()
	t := elem.Type()
	if len(ev.GetEventType()) == 0 {
		ev.SetEventType(t.Name())
	}
	if len(ev.GetAggregateID()) == 0 {
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			value, isAggregate := field.Tag.Lookup("aggregateID")
			if isAggregate && value == "true" {
				var svalue string
				switch field.Type.Kind() {
				case reflect.String:
					svalue = elem.Field(i).String()
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					svalue = fmt.Sprintf("%d", elem.Field(i).Int())
				case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					svalue = fmt.Sprintf("%d", elem.Field(i).Uint())
				default:
					return fmt.Errorf("aggregateID field must be of type string, int or uint %w", ErrInvalidEvent)
				}
				ev.SetAggregateID(domain + "-" + svalue)
				break
			}
		}
		ev.SetAggregateHash()
	}
	err = ev.Validate()
	return
}

// CommandToCommandRecord converts a command to a command record.
func CommandToCommandRecord(domain string, ev ICommand) (CommandRecord, error) {
	err := prepareCommand(domain, ev)
	if err != nil {
		return CommandRecord{}, err
	}
	data, err := json.Marshal(ev)
	if err != nil {
		return CommandRecord{}, err
	}
	return CommandRecord{
		RecordBase: RecordBase{
			ID:          ev.GetID(),
			AggregateID: ev.GetAggregateID(),
			EventType:   ev.GetEventType(),
			Data:        data,
			CreatedAt:   time.Now().UTC(),
		},
		AggregateHash: ev.GetAggregateHash(),
	}, nil
}

type CommandRequest struct {
	Name    string          `json:"name" validate:"required,gte=1,lte=100"`
	Payload json.RawMessage `json:"payload"`
}

func ParseCommandRequest(registry *Registry, r io.Reader) (ICommand, error) {
	var req CommandRequest
	if err := json.NewDecoder(r).Decode(&req); err != nil {
		return nil, fmt.Errorf("%w %s", ErrInvalidCommand, err.Error())
	}
	if err := core.Validate(req); err != nil {
		return nil, fmt.Errorf("%w %s", ErrInvalidCommand, err.Error())
	}
	conv, ok := registry.GetCommand(req.Name)
	if !ok {
		return nil, fmt.Errorf("%w %s", ErrInvalidCommand, "command not found")
	}
	return conv(req.Payload)
}
