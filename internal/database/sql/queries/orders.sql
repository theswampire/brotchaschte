-- name: GetOrder :one
SELECT * FROM raw_orders
WHERE id = ? LIMIT 1;

-- name: GetOrders :many
SELECT * FROM raw_orders 
LEFT JOIN geocode_errors 
ON raw_orders.id = geocode_errors.order_id;

-- name: CreateOrder :exec
INSERT INTO raw_orders (id, lastName, firstName, phoneNumber, email, amountSizeS, amountSizeM, amountSizeL, address, city, postalCode, deliveryMethod, remarks, orderMethod, vendorGroup)
VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?);

-- name: GetOrderStats :one
SELECT *, total_s + total_m + total_l AS total_amount FROM (
SELECT 
    COUNT(*) AS total_orders,
    SUM(amountSizeS) AS total_s, 
    SUM(amountSizeM) AS total_m,
    SUM(amountSizeL) AS total_l
FROM raw_orders) AS agg;

-- name: CreateOrderIfNotExist :exec
INSERT INTO raw_orders (id, lastName, firstName, phoneNumber, email, amountSizeS, amountSizeM, amountSizeL, address, city, postalCode, deliveryMethod, remarks, orderMethod, vendorGroup) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)
    ON CONFLICT(id)
    DO NOTHING;
