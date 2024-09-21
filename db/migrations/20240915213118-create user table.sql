
-- +migrate Up
CREATE TABLE user (
    id INTEGER PRIMARY KEY,
    userName TEXT,
    passwordHash TEXT,
    isActive BOOLEAN,
    isAdmin BOOLEAN);
-- +migrate Down
DROP TABLE user;