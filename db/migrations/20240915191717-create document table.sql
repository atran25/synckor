-- +migrate Up
CREATE TABLE document (
    hash TEXT,
    progress TEXT,
    percentage NUMERIC,
    device TEXT,
    deviceID TEXT,
    timestamp DATETIME
);

-- +migrate Down
DROP TABLE document;