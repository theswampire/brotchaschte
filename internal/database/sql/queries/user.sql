-- name: GetUser :one
SELECT * FROM user
WHERE email = ? LIMIT 1;

-- name: GetUsers :many
SELECT * FROM user;

-- name: CreateUser :exec
INSERT INTO user (email, password_hash)
VALUES (?,?);
