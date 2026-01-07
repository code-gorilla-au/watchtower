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
RETURNING *;

-- name: AddProductToOrganisation :exec
INSERT INTO product_organisations (product_id,
                                   organisation_id)
VALUES (?,
        ?);

-- name: DeleteProductOrganisationByOrgID :exec
DELETE
FROM product_organisations
WHERE organisation_id = ?;

-- name: UpdateProduct :one
UPDATE products
SET name        = ?,
    tags        = ?,
    description = ?,
    updated_at  = CAST(strftime('%s', 'now') AS INTEGER)
WHERE id = ?
RETURNING *;

-- name: UpdateProductSync :exec
UPDATE products
SET updated_at = CAST(strftime('%s', 'now') AS INTEGER)
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
SELECT o.*
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
RETURNING *;


-- name: UpdateRepo :one
UPDATE repositories
SET name       = ?,
    url        = ?,
    topic      = ?,
    owner      = ?,
    updated_at = CAST(strftime('%s', 'now') AS INTEGER)
WHERE id = ?
RETURNING *;

-- name: GetRepoByName :one
SELECT *
FROM repositories
WHERE name = ?
LIMIT 1;

-- name: GetReposByProductID :many
SELECT r.*, p.name as product_name
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
        ?,
        CAST(strftime('%s', 'now') AS INTEGER))
RETURNING *;

-- name: UpdatePullRequest :one
UPDATE pull_requests
SET title           = ?,
    repository_name = ?,
    url             = ?,
    state           = ?,
    author          = ?,
    merged_at       = ?,
    updated_at      = CAST(strftime('%s', 'now') AS INTEGER)
WHERE id = ?
RETURNING *;

-- name: GetPullRequestByExternalID :one
SELECT *
FROM pull_requests
WHERE external_id = ?
LIMIT 1;

-- name: GetPullRequestByProductIDAndState :many
SELECT pr.*, r.topic as tag, p.name as product_name
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
SELECT pr.*, r.topic as tag, p.name as product_name
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

-- name: DeletePullRequestsByProductID :exec
DELETE
FROM pull_requests
WHERE external_id IN (SELECT pr.external_id
                      FROM pull_requests pr
                               JOIN repositories r ON r.name = pr.repository_name
                               JOIN products p ON p.id = ?
                          AND JSON_VALID(p.tags)
                          AND EXISTS (SELECT 1
                                      FROM JSON_EACH(p.tags)
                                      WHERE JSON_EACH.value = r.topic));

-- name: CreateSecurity :one
INSERT INTO securities (external_id,
                        repository_name,
                        package_name,
                        state, severity,
                        patched_version,
                        fixed_at,
                        created_at,
                        updated_at)
VALUES (?,
        ?,
        ?,
        ?,
        ?,
        ?,
        ?,
        ?,
        CAST(strftime('%s', 'now') AS INTEGER))
ON CONFLICT (external_id) DO UPDATE SET repository_name = excluded.repository_name,
                                        package_name    = excluded.package_name,
                                        state           = excluded.state,
                                        severity        = excluded.severity,
                                        patched_version = excluded.patched_version,
                                        fixed_at        = excluded.fixed_at,
                                        updated_at      = CAST(strftime('%s', 'now') AS INTEGER)
RETURNING *;


-- name: GetSecurityByProductIDAndState :many
SELECT s.*, r.topic as tag, p.name as product_name
FROM securities s
         JOIN repositories r ON r.name = s.repository_name
         JOIN products p ON p.id = ?
    AND JSON_VALID(p.tags)
    AND EXISTS (SELECT 1
                FROM JSON_EACH(p.tags)
                WHERE JSON_EACH.value = r.topic)
WHERE s.state = ?
ORDER BY s.created_at DESC;

-- name: GetSecurityByOrganisationAndState :many
SELECT s.*, r.topic as tag, p.name as product_name
FROM securities s
         JOIN repositories r ON r.name = s.repository_name
         JOIN product_organisations po
         JOIN products p ON p.id = po.product_id
    AND JSON_VALID(p.tags)
    AND EXISTS (SELECT 1
                FROM JSON_EACH(p.tags)
                WHERE JSON_EACH.value = r.topic)
WHERE po.organisation_id = ?
  AND s.state = ?
ORDER BY s.created_at DESC;

-- name: DeleteSecurityByProductID :exec
DELETE
FROM securities
WHERE external_id IN (SELECT s.external_id
                      FROM securities s
                               JOIN repositories r ON r.name = s.repository_name
                               JOIN products p ON p.id = ?
                          AND JSON_VALID(p.tags)
                          AND EXISTS (SELECT 1
                                      FROM JSON_EACH(p.tags)
                                      WHERE JSON_EACH.value = r.topic));