
-- +migrate Up
CREATE TEMPORARY TABLE temp AS
    SELECT id,
           userName,
           passwordHash,
           isActive,
           isAdmin
    FROM user;

DROP TABLE user;

CREATE TABLE user (
    id INTEGER PRIMARY KEY,
    userName TEXT NOT NULL,
    passwordHash TEXT NOT NULL,
    isActive BOOLEAN NOT NULL,
    isAdmin BOOLEAN NOT NULL
);

INSERT INTO user (
    id,
    userName,
    passwordHash,
    isActive,
    isAdmin
) SELECT
    id,
    userName,
    passwordHash,
    isActive,
    isAdmin
FROM temp;

DROP TABLE temp;

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
    hash TEXT NOT NULL,
    progress TEXT NOT NULL,
    percentage NUMERIC NOT NULL,
    device TEXT NOT NULL,
    deviceID TEXT NOT NULL,
    timestamp DATETIME NOT NULL,
    userID INTEGER NOT NULL,
    FOREIGN KEY (userID) REFERENCES user (id) ON DELETE CASCADE
);

INSERT INTO document (
    hash,
    progress,
    percentage,
    device,
    deviceID,
    timestamp,
    userID
) SELECT
    hash,
    progress,
    percentage,
    device,
    deviceID,
    timestamp,
    userID
FROM temp;

DROP TABLE temp;

-- +migrate Down
CREATE TEMPORARY TABLE temp AS
    SELECT id,
           userName,
           passwordHash,
           isActive,
           isAdmin
    FROM user;

DROP TABLE user;

CREATE TABLE user (
    id INTEGER PRIMARY KEY,
    userName TEXT,
    passwordHash TEXT,
    isActive BOOLEAN,
    isAdmin BOOLEAN
);

INSERT INTO user (
    id,
    userName,
    passwordHash,
    isActive,
    isAdmin
) SELECT
    id,
    userName,
    passwordHash,
    isActive,
    isAdmin
FROM temp;

DROP TABLE temp;

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
    timestamp DATETIME,
    userID INTEGER
);

INSERT INTO document (
    hash,
    progress,
    percentage,
    device,
    deviceID,
    timestamp,
    userID
) SELECT
    hash,
    progress,
    percentage,
    device,
    deviceID,
    timestamp,
    userID
FROM temp;

DROP TABLE temp;