package es

import (
	"encoding/json"
	"time"
)

type BusMessage struct {
	Key       []byte
	Data      []byte
	Timestamp time.Time
}

func CommandRecordToBusMessage(cr CommandRecord) (BusMessage, error) {
	data, err := json.Marshal(cr)
	if err != nil {
		return BusMessage{}, err
	}
	return BusMessage{
		Key:  []byte(cr.AggregateID),
		Data: data,
	}, nil
}

func BusMessageToCommandRecord(msg BusMessage, cr *CommandRecord) error {
	err := json.Unmarshal(msg.Data, cr)
	if err != nil {
		return err
	}
	switch {
	case msg.Timestamp.IsZero():
		cr.CreatedAt = time.Now().UTC()
	default:
		cr.CreatedAt = msg.Timestamp
	}
	return nil
}

func BusMessagesToCommandRecords(msgs []BusMessage) ([]CommandRecord, error) {
	crs := make([]CommandRecord, len(msgs))
	for i := range msgs {
		err := BusMessageToCommandRecord(msgs[i], &crs[i])
		if err != nil {
			return nil, err
		}
	}
	return crs, nil
}
