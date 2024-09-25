
-- +migrate Up
CREATE TABLE user_account (
    username TEXT PRIMARY KEY ,
    password_hash TEXT NOT NULL
);
CREATE TABLE document_information (
    hash TEXT NOT NULL,
    progress TEXT NOT NULL,
    percentage NUMERIC NOT NULL,
    device TEXT NOT NULL,
    device_id TEXT NOT NULL,
    timestamp DATETIME NOT NULL,
    username TEXT NOT NULL,
    FOREIGN KEY (username) REFERENCES user_account (username) ON DELETE CASCADE,
    PRIMARY KEY (hash, username)
);
-- +migrate Down
DROP TABLE user_account;
DROP TABLE document_information;