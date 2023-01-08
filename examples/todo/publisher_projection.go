package todo

import (
	"context"
	"database/sql"
	"time"

	"github.com/gosom/kit/es"
	"github.com/gosom/kit/logging"
	"github.com/gosom/kit/sqldb"
)

type ProjectionBuilder struct {
	db       *sqldb.DB
	registry *es.Registry

	log logging.Logger
}

func NewProjectionBuilder(db *sqldb.DB, registry *es.Registry) *ProjectionBuilder {
	return &ProjectionBuilder{
		db:       db,
		registry: registry,
		log: logging.Get().With("component", "todo_projection").Level(
			logging.DEBUG,
		),
	}
}

func (p *ProjectionBuilder) Publish(ctx context.Context, records ...es.EventRecord) error {
	events, err := es.EventRecordsToEvents(p.registry, records)
	if err != nil {
		return err
	}
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()
	for i := range events {
		switch e := events[i].(type) {
		case *TodoCreated:
			if err := p.processTodoCreated(ctx, tx, records[i].CreatedAt, e); err != nil {
				return err
			}
		case *TodoStatusUpdated:
			if err := p.processTodoStatusUpdated(ctx, tx, records[i].CreatedAt, e); err != nil {
				return err
			}
		default:
			p.log.Warn("unknown event", "event", e)
		}
	}
	return tx.Commit()
}

func (p *ProjectionBuilder) Name() string {
	return "todo_projection"
}

func (p *ProjectionBuilder) processTodoCreated(ctx context.Context, tx *sql.Tx, ts time.Time, e *TodoCreated) error {
	const q = `INSERT INTO todos
	(id, title, status, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5)`
	_, err := tx.ExecContext(ctx, q, e.ID, e.Title, "open", ts, ts)
	return err
}

func (p *ProjectionBuilder) processTodoStatusUpdated(ctx context.Context, tx *sql.Tx, ts time.Time, e *TodoStatusUpdated) error {
	const q = `UPDATE todos
	SET status = $1, updated_at = $2
	WHERE id = $3`
	_, err := tx.ExecContext(ctx, q, "completed", ts, e.ID)
	return err
}
