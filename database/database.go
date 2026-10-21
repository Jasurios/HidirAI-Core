package database

import (
	"database/sql"
	_ "modernc.org/sqlite"
)

func Open(namedb string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", "file:./database/"+namedb+"?_pragma=journal_mode(WAL)")
	if err != nil {
		return nil, err
	}

	return db, nil
}

func Init(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS keys (
			type TEXT NOT NULL,
			user TEXT NOT NULL,
			key TEXT NOT NULL
		)
	`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS blacklist (
			user TEXT NOT NULL
		)
	`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_keys_key ON keys(key);
		CREATE UNIQUE INDEX IF NOT EXISTS idx_keys_user ON keys(user);
		CREATE UNIQUE INDEX IF NOT EXISTS idx_blacklist_user ON blacklist(user);
	`)
	if err != nil {
		return err
	}

	return nil
}
