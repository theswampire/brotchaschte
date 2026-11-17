-- name: GetFailedGeocodes :many
SELECT * FROM geocode_errors;

-- name: GetFailedGeocodeOrders :many
SELECT * FROM geocode_errors
JOIN raw_orders ON geocode_errors.order_id = raw_orders.id;

-- name: GetGeocodeWithSuggestionOrder :one
SELECT * FROM raw_orders
LEFT JOIN geocode_errors ON raw_orders.id = geocode_errors.order_id
LEFT JOIN order_coordinates ON raw_orders.id = order_coordinates.order_id
WHERE raw_orders.id = ?;

-- name: CleanUpErrors :exec
DELETE FROM geocode_errors 
WHERE order_id IN (SELECT order_id FROM order_coordinates);

-- name: GetNotGeocodedOrdersIncludeFailed :many
SELECT * FROM raw_orders 
LEFT JOIN geocode_errors ON raw_orders.id = geocode_errors.order_id
WHERE raw_orders.id NOT IN (SELECT order_id AS id FROM order_coordinates);

-- name: GetNotGeocodedOrdersExcludeFailed :many
SELECT * FROM raw_orders
WHERE id NOT IN (
    SELECT order_id AS id FROM order_coordinates
    UNION
    SELECT order_id AS id  FROM geocode_errors
);

-- name: GetGeocodedOrders :many
SELECT * FROM raw_orders
JOIN order_coordinates ON raw_orders.id = order_coordinates.order_id;

-- name: CreateGeocodeError :exec
INSERT INTO geocode_errors(order_id, message) VALUES (?1, ?2)
    ON CONFLICT(order_id)
    DO 
        UPDATE SET message = ?2
        WHERE order_id = ?1;

-- name: CreateGeocodeErrorWithSuggestion :exec
INSERT INTO geocode_errors(order_id, message, suggested_northing, suggested_easting, suggested_label) VALUES (?1, ?2, ?3, ?4, ?5)
    ON CONFLICT(order_id)
    DO 
        UPDATE SET message = ?2, suggested_northing = ?3, suggested_easting = ?4, suggested_label = ?5
        WHERE order_id = ?1;

-- name: DeleteGeocodeError :exec
DELETE FROM geocode_errors
WHERE order_id = ?;

-- name: CreateGeocode :exec
INSERT INTO order_coordinates(order_id, x, y, label) VALUES (?, ?, ?, ?);

-- name: UpsertGeocode :exec
INSERT INTO order_coordinates(order_id, x, y, label) VALUES (?1, ?2, ?3, ?4)
    ON CONFLICT(order_id)
    DO 
        UPDATE SET x = ?2, y = ?3, label = ?4
        WHERE order_id = ?1;

-- name: ManualGeocode :one
INSERT INTO order_coordinates(order_id, x, y, label)
SELECT ?1, ?2, ?3, address || ' ' || postalCode || ' ' || city
FROM raw_orders
WHERE id = ?1
ON CONFLICT(order_id)
DO 
    UPDATE SET x = excluded.x, y = excluded.y, label = excluded.label
    WHERE order_id = ?1
RETURNING *;

-- name: UpdateAddress :one
UPDATE raw_orders SET address = ?, postalCode = ?, city = ?, remarks = ?
WHERE id = ?
RETURNING *;
