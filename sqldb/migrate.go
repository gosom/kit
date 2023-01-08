package sqldb

import (
	"context"
	"embed"

	"github.com/golang-migrate/migrate/v4"
	mpostgres "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

func Migrate(ctx context.Context, db *DB, migrationsTable string, fs embed.FS) error {
	config := mpostgres.Config{}
	if migrationsTable != "" {
		config.MigrationsTable = migrationsTable
	}
	driver, err := mpostgres.WithInstance(db.Conn(), &config)
	if err != nil {
		return err
	}
	source, err := iofs.New(fs, "migrations")
	if err != nil {
		return err
	}
	m, err := migrate.NewWithInstance("iofs", source, "postgres", driver)
	if err != nil {
		return err
	}
	err = m.Up()
	switch err {
	case nil, migrate.ErrNoChange:
		return nil
	default:
		return err
	}
}
