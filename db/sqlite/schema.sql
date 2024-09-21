CREATE TABLE IF NOT EXISTS "gorp_migrations" ("id" varchar(255) not null primary key, "applied_at" datetime);
CREATE TABLE user (
    id INTEGER PRIMARY KEY,
    userName TEXT NOT NULL,
    passwordHash TEXT NOT NULL,
    isActive BOOLEAN NOT NULL,
    isAdmin BOOLEAN NOT NULL
);
CREATE TABLE document (
    hash TEXT NOT NULL,
    progress TEXT NOT NULL,
    percentage NUMERIC NOT NULL,
    device TEXT NOT NULL,
    deviceID TEXT NOT NULL,
    timestamp DATETIME NOT NULL,
    userID INTEGER NOT NULL,
    FOREIGN KEY (userID) REFERENCES user (id) ON DELETE CASCADE
);
