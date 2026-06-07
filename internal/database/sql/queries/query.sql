-- name: GetRandomAddresses :many
SELECT * FROM avz
ORDER BY random()
LIMIT ?;
