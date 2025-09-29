package internal

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// DatabaseFile is the default filename for the SQLite database used by the
// server.
const DatabaseFile = "quego.db"

// MigrateDatabase applies database schema migrations.
func MigrateDatabase() error {
	migrations, err := migrate.New(
		"file://migrations", // where the migration files are located
		fmt.Sprintf("sqlite3://%s", DatabaseFile),
	)
	if err != nil {
		return fmt.Errorf("failed to initialize migrations: %w", err)
	}
	if err := migrations.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}
	return nil
}
