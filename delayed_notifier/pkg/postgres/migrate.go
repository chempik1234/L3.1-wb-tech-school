package postgres

import (
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // for postgres driver
	_ "github.com/golang-migrate/migrate/v4/source/file"       // for file driver (search for sql)
)

// MigrateUp runs all migrations by path and DSN
//
//	connString = DSN
//	sourcePath = "file:///app/db/migrations"
func MigrateUp(connString string, sourcePath string) error {
	m, err := migrate.New(
		sourcePath,
		connString,
	)
	if err != nil {
		return fmt.Errorf("unable to create migrations table: %v", err)
	}
	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("unable to apply migrations: %v", err)
	}
	return nil
}
