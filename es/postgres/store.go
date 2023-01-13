package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/gosom/kit/es"
	"github.com/gosom/kit/es/assets"
	"github.com/gosom/kit/logging"
	"github.com/gosom/kit/sqldb"
)

type aggregateVersion struct {
	AggregateID string
	Version     int
}

func (o *aggregateVersion) Bind() []any {
	return []any{&o.AggregateID, &o.Version}
}

type commandRecord struct {
	es.CommandRecord
	Rn int
}

func (o *commandRecord) Bind() []any {
	ans := o.CommandRecord.Bind()
	ans = append(ans, &o.Rn)
	return ans
}

var _ es.EventStore = (*EventStore)(nil)

type EventStore struct {
	db  *sqldb.DB
	log logging.Logger
}

func NewEventStore(db *sqldb.DB) *EventStore {
	return &EventStore{db: db, log: logging.Get().With("component", "store")}
}

func (e *EventStore) SaveCommandRecords(ctx context.Context, records ...es.CommandRecord) ([]string, error) {
	valueStrings := make([]string, 0, len(records))
	valueArgs := make([]interface{}, 0, len(records)*6)
	for i := range records {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d)",
			i*6+1, i*6+2, i*6+3, i*6+4, i*6+5, i*6+6))
		valueArgs = append(valueArgs,
			records[i].ID,
			records[i].AggregateID,
			records[i].EventType,
			records[i].Data,
			records[i].CreatedAt,
			records[i].AggregateHash)
	}
	stmt := fmt.Sprintf(saveCommandsStmt, strings.Join(valueStrings, ","))
	rows, err := e.db.Conn().QueryContext(ctx, stmt, valueArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

func (e *EventStore) SaveCommand(ctx context.Context, domain string, cmd es.ICommand) (string, error) {
	rec, err := es.CommandToCommandRecord(domain, cmd)
	if err != nil {
		return "", err
	}
	ids, err := e.SaveCommandRecords(ctx, rec)
	if err != nil {
		return "", err
	}
	if len(ids) == 0 {
		return "", fmt.Errorf("no command records saved")
	}
	return ids[0], nil
}

func (e *EventStore) GetCommand(ctx context.Context, commandID string) (es.CommandRecord, error) {
	record, err := sqldb.QueryRow[es.CommandRecord](ctx, e.db.Conn(), getCommandStmt, commandID)
	return record, err
}

func (e *EventStore) Migrate(ctx context.Context) error {
	return sqldb.Migrate(ctx, e.db, "es_schema_migrations", assets.Migrations)
}

func (e *EventStore) SelectForProcessing(ctx context.Context, workers int, limit int) ([][]es.CommandRecord, error) {
	ans := make([][]es.CommandRecord, workers)
	for i := 0; i < workers; i++ {
		ans[i] = make([]es.CommandRecord, 0, limit)
	}
	records, err := sqldb.Query[commandRecord](ctx, e.db.Conn(), selectCommandsToProcess, workers, limit)
	if err != nil {
		return nil, err
	}
	if len(records) == 0 {
		return ans, nil
	}
	current := 0
	ans[0] = append(ans[0], records[0].CommandRecord)
	for i := 1; i < len(records); i++ {
		if records[i].Rn <= current {
			current++
		}
		ans[current] = append(ans[current], records[i].CommandRecord)
	}
	return ans, nil
}

func (e *EventStore) StoreCommandResults(ctx context.Context, commandID string, expectedVersion int, events ...es.EventRecord) error {
	tx, err := e.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()
	if len(events) > 0 {
		rs, err := tx.ExecContext(ctx, checkVersionStmt, len(events), events[0].AggregateID, expectedVersion)
		if err != nil {
			return fmt.Errorf("error updating version: %w", err)
		}
		affected, err := rs.RowsAffected()
		if err != nil {
			return err
		}
		if affected == 0 {
			return es.ErrWrongExpectedVersion
		}
	}
	for i := range events {
		if _, err := tx.ExecContext(ctx, saveEventsStmt, events[i].ID, commandID, events[i].AggregateID, events[i].Version, events[i].EventType, events[i].Data); err != nil {
			return fmt.Errorf("Error saving event %s: %w", events[i].ID, err)
		}
	}
	if _, err := tx.ExecContext(ctx, updateCommandStatusStmt, "finished", commandID); err != nil {
		return fmt.Errorf("error updating commandStatus: %w", err)
	}
	return tx.Commit()
}

func (e *EventStore) GetOrCreateVersion(ctx context.Context, aggregateID string) (int, error) {
	rec, err := sqldb.QueryRow[aggregateVersion](ctx, e.db.Conn(), getOrCreateAggregateVersionStmt, aggregateID)
	if err != nil {
		return 0, err
	}
	return rec.Version, nil
}

func (e *EventStore) InsertSubscription(ctx context.Context, subscription string) (es.Subscription, error) {
	sub, err := sqldb.QueryRow[es.Subscription](ctx, e.db.Conn(), insertSubStmt, subscription)
	return sub, err
}

func (e *EventStore) SelectEventsForSubscription(ctx context.Context, subscription es.Subscription, limit int) ([]es.EventRecord, error) {
	records, err := sqldb.Query[es.EventRecord](ctx, e.db.Conn(), selectEventsForSubStmt, subscription.Group, limit)
	return records, err
}

func (e *EventStore) UpdateSubscription(ctx context.Context, group string, lastSeen string) (es.Subscription, error) {
	sub, err := sqldb.QueryRow[es.Subscription](ctx, e.db.Conn(), updateSubStmt, group, lastSeen)
	return sub, err
}

func (e *EventStore) LoadEvents(ctx context.Context, aggregateID string) ([]es.EventRecord, error) {
	records, err := sqldb.Query[es.EventRecord](ctx, e.db.Conn(), loadEventsStmt, aggregateID)
	return records, err
}
