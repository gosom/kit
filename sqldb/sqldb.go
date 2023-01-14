package sqldb

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

// Bindable is an interface that can be used to bind a struct to a sql query.
type Bindable[T any] interface {
	*T
	Bind() []any
}

// RowScanner is an interface that can be used to scan a row into a struct.
type RowScanner interface {
	Scan(dest ...interface{}) error
}

type DBTX interface {
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
}

// DB is a wrapper around sql.DB that provides some additional functionality.
type DB struct {
	pool *sql.DB

	ctx    context.Context
	cancel func()

	DriverName string
	DSN        string

	Now func() time.Time
}

// NewDB returns a new DB.
func NewDB(driver, dsn string) *DB {
	ans := DB{
		DriverName: driver,
		DSN:        dsn,
		Now: func() time.Time {
			return time.Now().UTC()
		},
	}
	ans.ctx, ans.cancel = context.WithCancel(context.Background())
	return &ans
}

// SetPool sets the database pool.
func (o *DB) SetPool(pool *sql.DB) {
	o.pool = pool
}

// Open opens the database connection.
func (o *DB) Open() (err error) {
	if o.DriverName == "" {
		return errors.New("DriverName is empty")
	}
	if o.DSN == "" {
		return errors.New("DSN is empty")
	}
	o.pool, err = sql.Open(o.DriverName, o.DSN)
	if err != nil {
		return
	}
	err = o.pool.Ping()
	return
}

func (o *DB) Conn() *sql.DB {
	return o.pool
}

// Close closes the database connection.
func (o *DB) Close() error {
	o.cancel()
	if o.pool == nil {
		return nil
	}
	return o.pool.Close()
}

// BeginTx begins a transaction.
func (o *DB) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	return o.pool.BeginTx(ctx, opts)
}

// Query returns a slice of items from the database.
func Query[T any, PT Bindable[T]](ctx context.Context, tx DBTX, q string, args ...any) ([]T, error) {
	rows, err := tx.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []T
	for rows.Next() {
		var item T
		var pt PT = &item
		if err := rows.Scan(pt.Bind()...); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

// QueryRow returns a single item from the database.
func QueryRow[T any, PT Bindable[T]](ctx context.Context, tx DBTX, q string, args ...any) (T, error) {
	var item T
	var pt PT = &item
	row := tx.QueryRowContext(ctx, q, args...)
	if err := row.Scan(pt.Bind()...); err != nil {
		return item, err
	}
	return item, nil
}
