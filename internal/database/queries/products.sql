-- name: CreateProduct :one
INSERT INTO products (name,
                      tags,
                      description,
                      created_at,
                      updated_at)
VALUES (?,
        ?,
        ?,
        CAST(strftime('%s', 'now') AS INTEGER),
        CAST(strftime('%s', 'now') AS INTEGER))
ON CONFLICT (name) DO UPDATE SET name       = excluded.name,
                                 tags       = excluded.tags,
                                 updated_at = CAST(strftime('%s', 'now') AS INTEGER)
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
DELETE
FROM products
where id = ?;

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
        CAST(strftime('%s', 'now') AS INTEGER))
ON CONFLICT (name) DO UPDATE SET name       = excluded.name,
                                 url        = excluded.url,
                                 topic      = excluded.topic,
                                 owner      = excluded.owner,
                                 updated_at = CAST(strftime('%s', 'now') AS INTEGER)
RETURNING *;


-- name: GetReposByProductID :many
SELECT r.*
FROM repositories r
         JOIN products p ON p.id = ?
    AND JSON_VALID(p.tags)
    AND EXISTS (SELECT 1
                FROM JSON_EACH(p.tags)
                WHERE JSON_EACH.value = r.topic);


-- name: DeleteReposByProductID :exec
DELETE
FROM repositories
WHERE topic IN (SELECT JSON_EACH.value
                FROM products p, JSON_EACH(p.tags)
                WHERE p.id = ?
                  AND JSON_VALID(p.tags));

-- name: CreatePullRequest :one
INSERT INTO pull_requests (external_id,
                           title,
                           repository_name,
                           url,
                           state,
                           author,
                           merged_at,
                           created_at,
                           updated_at)
VALUES (?,
        ?,
        ?,
        ?,
        ?,
        ?,
        ?,
        CAST(strftime('%s', 'now') AS INTEGER),
        CAST(strftime('%s', 'now') AS INTEGER))
ON CONFLICT (external_id) DO UPDATE SET title           = excluded.title,
                                        repository_name = excluded.repository_name,
                                        url             = excluded.url,
                                        state           = excluded.state,
                                        author          = excluded.author,
                                        merged_at       = excluded.merged_at,
                                        updated_at      = CAST(strftime('%s', 'now') AS INTEGER)
RETURNING *;

-- name: GetPullRequestByProductIDAndState :many
SELECT pr.*
FROM pull_requests pr
         JOIN repositories r ON r.name = pr.repository_name
         JOIN products p ON p.id = ?
    AND JSON_VALID(p.tags)
    AND EXISTS (SELECT 1
                FROM JSON_EACH(p.tags)
                WHERE JSON_EACH.value = r.topic)
WHERE pr.state = ?
ORDER BY pr.created_at DESC;

-- name: GetPullRequestsByOrganisationAndState :many
SELECT pr.*
FROM pull_requests pr
         JOIN repositories r ON r.name = pr.repository_name
         JOIN product_organisations po
         JOIN products p ON p.id = po.product_id
    AND JSON_VALID(p.tags)
    AND EXISTS (SELECT 1
                FROM JSON_EACH(p.tags)
                WHERE JSON_EACH.value = r.topic)
WHERE po.organisation_id = ?
AND pr.state = ?
ORDER BY pr.created_at DESC;