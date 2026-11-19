-- name: GetUser :one
SELECT * FROM user
WHERE email = ? LIMIT 1;

-- name: GetUsers :many
SELECT * FROM user;

-- name: CreateUser :exec
INSERT INTO user (email, password_hash)
VALUES (?,?);

-- name: DeleteUser :exec
DELETE FROM user WHERE email = ?;

-- name: GetUserCount :one
SELECT COUNT(*) FROM user;

-- name: UpdatePassword :exec
UPDATE user SET password_hash = ?
WHERE email = ?;
