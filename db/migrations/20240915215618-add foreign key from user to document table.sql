
-- +migrate Up
CREATE TEMPORARY TABLE temp AS
    SELECT hash,
           progress,
           percentage,
           device,
           deviceID,
           timestamp
    FROM document;

DROP TABLE document;

CREATE TABLE document (
    hash TEXT,
    progress TEXT,
    percentage NUMERIC,
    device TEXT,
    deviceID TEXT,
    timestamp DATETIME,
    userID INTEGER,
    FOREIGN KEY (userID) REFERENCES user (id)
);

DROP TABLE temp;

-- +migrate Down
CREATE TEMPORARY TABLE temp AS
    SELECT hash,
           progress,
           percentage,
           device,
           deviceID,
           timestamp,
           userID
    FROM document;

DROP TABLE document;

CREATE TABLE document (
    hash TEXT,
    progress TEXT,
    percentage NUMERIC,
    device TEXT,
    deviceID TEXT,
    timestamp DATETIME
);

DROP TABLE temp;