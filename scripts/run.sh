#!/bin/bash
set -e

# Restore the database if it does not already exist.
if [ -f /app/db/sqlite/db.sqlite ]; then
	echo "Database already exists, skipping restore"
else
	echo "No database found, restoring from replica if exists"
	litestream restore -v -if-replica-exists -o /app/db/sqlite/db.sqlite "${REPLICA_URL}"

	# Check if the restore command actually created the database file.
      if [ -f /app/db/sqlite/db.sqlite ]; then
          echo "Database restored from replica"
      else
          echo "No replica found or restore did not create the database, creating an empty SQLite database"
          sqlite3 /app/db/sqlite/db.sqlite "VACUUM"
      fi
fi

# Run the migrations
echo "Running migrations"
/app/migrate

# Run litestream with your app as the subprocess.
echo "Starting litestream"
exec litestream replicate -exec "/app/synckor"
