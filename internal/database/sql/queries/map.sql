-- name: GetOrderCoordinates :many
SELECT r.*, oc.x AS x, oc.y AS y, oc.label AS address_label, d.delivery_id AS delivery_id
FROM raw_orders AS r
JOIN order_coordinates AS oc  ON r.id = oc.order_id
LEFT JOIN delivery_orders AS d ON r.id = d.order_id
ORDER BY oc.x DESC;

-- name: GetDelivery :one
SELECT * FROM delivery
WHERE id = ?
LIMIT 1;

-- name: GetDeliveries :many
SELECT * FROM delivery
ORDER BY id DESC;


-- name: CreateDelivery :one
INSERT INTO delivery(label, until, color) VALUES (?, ?, ?)
RETURNING *;

-- name: UpdateDelivery :one
UPDATE delivery
SET until = ?, label = ?
WHERE id = ?
RETURNING *;

-- name: UpdateDeliveryColor :one
UPDATE delivery
SET color = ?
WHERE id = ?
RETURNING *;

-- name: DeleteDeliveryOrderDeselected :exec
DELETE FROM delivery_orders
WHERE order_id = ? AND delivery_id = ?;

-- name: CreateDeliveryOrderSelected :exec
INSERT INTO delivery_orders(order_id, delivery_id) VALUES (?1, ?2)
    ON CONFLICT(order_id)
    DO 
        UPDATE SET delivery_id = ?2
        WHERE order_id = ?1;
