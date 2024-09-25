-- name: GetUser :one
SELECT * FROM user_account WHERE username = ?;

-- name: GetUserWithPassword :one
SELECT * FROM user_account WHERE username = ? AND password_hash = ?;

-- name: CreateUser :one
INSERT INTO user_account (username, password_hash) VALUES (?, ?) RETURNING *;

-- name: GetDocument :one
SELECT * FROM document_information WHERE hash = ? and username = ?;

-- name: CreateDocument :one
INSERT INTO document_information (hash, progress, percentage, device, device_id, timestamp, username) VALUES (?, ?, ?, ?, ?, ?, ?) RETURNING *;

-- name: UpdateDocument :one
UPDATE document_information SET progress = ?, percentage = ?, device = ?, device_id = ?, timestamp = ? WHERE hash = ? and username = ? RETURNING *;