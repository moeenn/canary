package database

import (
	"canary/internal/database/migrations"
	"fmt"
	"os"
	"path"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

func Connect(dir, dbFilename string) (*sqlx.DB, error) {
	if err := os.MkdirAll(dir, 0750); err != nil {
		return nil, fmt.Errorf("failed to create database dir (%s): %w", dir, err)
	}

	dbFile := path.Join(dir, dbFilename)
	db, err := sqlx.Open("sqlite", dbFile)
	if err != nil {
		return nil, fmt.Errorf("failed to connect db: %w", err)
	}

	if err := migrations.Run(db); err != nil {
		return nil, err
	}

	return db, nil
}
