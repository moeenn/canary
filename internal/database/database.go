package database

import (
	"canary/internal/database/migrations"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

func Connect(dbFile string) (*sqlx.DB, error) {
	db, err := sqlx.Open("sqlite", dbFile)
	if err != nil {
		return nil, fmt.Errorf("failed to connect db: %w", err)
	}

	if err := migrations.Run(db); err != nil {
		return nil, err
	}

	return db, nil
}
