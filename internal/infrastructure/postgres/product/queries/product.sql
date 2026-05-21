-- name: CreateProduct :one
INSERT INTO products (id, name, description, price, sku)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, name, description, price, sku, created_at, updated_at;

-- name: GetProductByID :one
SELECT id, name, description, price, sku, created_at, updated_at
FROM products
WHERE id = $1;

-- name: ListProducts :many
SELECT id, name, description, price, sku, created_at, updated_at
FROM products
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountProducts :one
SELECT COUNT(*) AS count
FROM products;

-- name: UpdateProduct :one
UPDATE products
SET name        = $2,
    description = $3,
    price       = $4,
    sku         = $5,
    updated_at  = NOW()
WHERE id = $1
RETURNING id, name, description, price, sku, created_at, updated_at;

-- name: DeleteProduct :exec
DELETE FROM products
WHERE id = $1;

-- name: ExistsProductBySKU :one
SELECT EXISTS(SELECT 1 FROM products WHERE sku = $1) AS exists;
