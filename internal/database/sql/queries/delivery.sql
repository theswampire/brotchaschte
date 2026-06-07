-- name: GetOrderedDeliveries :many
SELECT orders.*, delivery.id AS delivery_id, delivery.until AS delivery_until, delivery.label AS delivery_label, coords.label AS geocoded_address
FROM raw_orders AS orders
LEFT JOIN delivery_orders
    ON orders.id = delivery_orders.order_id
LEFT JOIN delivery
    ON delivery_orders.delivery_id = delivery.id
LEFT JOIN order_coordinates AS coords
    ON coords.order_id = orders.id
ORDER BY delivery.id, orders.address ASC;
