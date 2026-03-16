package cron

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

type SQLiteCronStore struct {
	db *sql.DB
}

func NewSQLiteCronStore(dbPath string) (*SQLiteCronStore, error) {
	db, err := sql.Open("sqlite", dbPath+"?_journal_mode=WAL")
	if err != nil {
		return nil, fmt.Errorf("open cron sqlite store: %w", err)
	}

	store := &SQLiteCronStore{db: db}
	if err := store.initialize(); err != nil {
		_ = db.Close()
		return nil, err
	}

	return store, nil
}

func (s *SQLiteCronStore) Close() error {
	if s == nil || s.db == nil {
		return nil
	}
	return s.db.Close()
}
