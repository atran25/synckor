package main

import (
	"database/sql"
	"fmt"
	migrate "github.com/rubenv/sql-migrate"
	"log/slog"
	_ "modernc.org/sqlite"
)

func main() {
	migrations := &migrate.FileMigrationSource{
		Dir: "db/migrations",
	}
	db, err := sql.Open("sqlite", "db/sqlite/db.sqlite")
	if err != nil {
		// Handle errors!
		slog.Error("Opening database:", err)
	}

	n, err := migrate.Exec(db, "sqlite3", migrations, migrate.Up)
	if err != nil {
		// Handle errors!
		slog.Error("Applying migrations:", err)
	}
	fmt.Printf("Applied %d migrations!\n", n)
}
