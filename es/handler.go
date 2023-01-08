package es

//type Worker struct {
//	Registry   *MessageToEventRegistry
//	EventStore IEventStore
//}
//
//func (w *Worker) Process(ctx context.Context, key []byte, msg []byte, timestamp time.Time) error {
//	var rec EventRecord
//	if err := json.Unmarshal(msg, &rec); err != nil {
//		return fmt.Errorf("%s %w", err, ErrSkipEvent)
//	}
//	convFn, ok := w.Registry.Get(rec.EventType)
//	if !ok {
//		return fmt.Errorf("no converter for event type %s %w", rec.EventType, ErrSkipEvent)
//	}
//	event, err := convFn(rec)
//	if err != nil {
//		return fmt.Errorf("rec: %s %s %w", rec.EventType, err, ErrSkipEvent)
//	}
//	txstore, err := w.EventStore.Begin(ctx)
//	if err != nil {
//		return fmt.Errorf("begin tx %w", err)
//	}
//	defer txstore.Rollback(ctx)
//	h := w
//	h.EventStore = txstore
//	events, err := event.Handle(h)
//	if err != nil {
//		return fmt.Errorf("handle event %s %w", err, ErrSkipEvent)
//	}
//	records := make([]EventRecord, len(events))
//	for i := range events {
//		records[i], err = eventToRecord(events[i])
//	}
//	return nil
//}
//
//func (w *Worker) saveEvents(ctx context.Context, txStore IEventStore, events []IEvent) (err error) {
//	records := make([]EventRecord, len(events))
//	for i := range events {
//		records[i], err = eventToRecord(events[i])
//	}
//	err = txStore.Save(ctx, records)
//	return
//}
