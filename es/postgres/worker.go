package postgres

//type worker struct {
//	log        logging.Logger
//	dispatcher *es.CommandDispatcher
//}
//
//func NewWorker(domain string, store es.EventStore, registry *es.Registry) *worker {
//	return &worker{
//		log:        logging.Get().With("component", "worker"),
//		dispatcher: es.NewCommandDispatcher("domain", store, es.NewRegistry()),
//	}
//}
//
//func (w *worker) Process(ctx context.Context, key []byte, value []byte, timestamp time.Time) error {
//	w.log.Info("Processing message", "key", string(key), "value", string(value), "timestamp", timestamp)
//	return nil
//}
