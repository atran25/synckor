package database

import (
	"database/sql"
	"log/slog"
	_ "modernc.org/sqlite"
)

func GetConnection() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "db/sqlite/db.sqlite")
	if err != nil {
		// Handle errors!
		slog.Error("Opening database:", "Error", err)
		return nil, err
	}

	// Litestream specific configuration settings for making sqlite usable: https://litestream.io/tips/
	_, err = db.Exec("PRAGMA busy_timeout = 5000;")
	if err != nil {
		slog.Error("Setting busy_timeout = 5000:", "Error", err)
		return nil, err
	}
	_, err = db.Exec("PRAGMA journal_mode = WAL;")
	if err != nil {
		slog.Error("Setting journal_mode = WAL:", "Error", err)
		return nil, err
	}
	_, err = db.Exec("PRAGMA synchronous = NORMAL;")
	if err != nil {
		slog.Error("Setting synchronous = NORMAL:", "Error", err)
		return nil, err
	}
	return db, err
}
