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
		slog.Error("Opening database:", err)
		return nil, err
	}
	return db, err
}
