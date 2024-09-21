-- name: GetUser :one
SELECT * FROM user WHERE username = ?;

-- name: CreateUser :one
INSERT INTO user (username, passwordHash, isActive, isAdmin) VALUES (?, ?, ?, ?) RETURNING *;