-- name: CreateProduct :one
INSERT INTO products (name,
                      tags,
                      created_at,
                      updated_at)
VALUES (?,
        ?,
        CAST(strftime('%s', 'now') AS INTEGER),
        CAST(strftime('%s', 'now') AS INTEGER))
RETURNING *;

-- name: AddProductToOrganisation :exec
INSERT INTO product_organisations (product_id,
                                   organisation_id)
VALUES (?,
        ?);

-- name: UpdateProduct :exec
UPDATE products
SET name       = ?,
    tags       = ?,
    updated_at = CAST(strftime('%s', 'now') AS INTEGER)
WHERE id = ?;

-- name: GetProductByID :one
SELECT *
FROM products
WHERE id = ?
LIMIT 1;


-- name: DeleteProduct :exec
DELETE FROM products where id = ?;

-- name: ListProductsByOrganisation :many
SELECT p.*
FROM products p
         JOIN product_organisations po ON po.product_id = p.id
WHERE po.organisation_id = ?
ORDER BY p.name;

-- name: GetOrganisationForProduct :one
SELECT *
FROM product_organisations po
         JOIN organisations o ON o.id = po.organisation_id
WHERE po.product_id = ?
LIMIT 1;

-- name: CreateRepo :one
INSERT INTO repositories (name,
                          url,
                          topic,
                          owner,
                          created_at,
                          updated_at)
VALUES (?,
        ?,
        ?,
        ?,
        CAST(strftime('%s', 'now') AS INTEGER),
        CAST(strftime('%s', 'now') AS INTEGER)) RETURNING *;


-- name: GetReposByProductID :many
SELECT r.*
FROM repositories r
JOIN products p ON p.id = ? 
    AND JSON_VALID(p.tags) 
    AND EXISTS (
        SELECT 1 
        FROM JSON_EACH(p.tags) 
        WHERE JSON_EACH.value = r.topic
    );