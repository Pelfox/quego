package internal

import (
	"embed"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

// DatabaseFile is the default filename for the SQLite database used by the
// server.
const DatabaseFile = "quego.db"

//go:embed ../migrations/*.sql
var migrationsFS embed.FS

// MigrateDatabase applies database schema migrations.
func MigrateDatabase() error {
	data, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		return err
	}
	migrations, err := migrate.NewWithSourceInstance(
		"iofs",
		data,
		fmt.Sprintf("sqlite3://%s", DatabaseFile),
	)
	if err != nil {
		return fmt.Errorf("failed to initialize migrations: %w", err)
	}
	if err := migrations.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}
	return nil
}
