package migrations

import (
	"fmt"
	"log/slog"

	"github.com/jmoiron/sqlx"
)

const enablePragmasMigration string = `
    PRAGMA foreign_keys = ON;
    PRAGMA journal_mode = WAL;
`

const createLogsTableMigration string = `
	create table if not exists logs (
		id varchar(26) primary key not null -- ulid
		, time text not null
		, level varchar (10) not null
		, payload text not null
	);
`

const createLogsTableIndicesMigration string = `
	create index log_level_idx on logs (level);
	create index log_time_idx on logs (time);
`

var migrations = []string{
	enablePragmasMigration,
	createLogsTableMigration,
	createLogsTableIndicesMigration,
}

func Run(db *sqlx.DB) error {
	slog.Info("running db migrations")
	for _, m := range migrations {
		if _, err := db.Exec(m, nil); err != nil {
			return fmt.Errorf("migrations failed: %w", err)
		}
	}

	return nil
}
