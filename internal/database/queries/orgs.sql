-- name: CreateOrganisation :one
INSERT INTO organisations (friendly_name,
                           namespace,
                           default_org,
                           token,
                           created_at,
                           updated_at)
VALUES (?,
        ?,
        true,
        ?,
        unixepoch('now'),
        unixepoch('now'))
RETURNING *;

-- name: ListOrganisations :many
SELECT *
FROM organisations
ORDER BY friendly_name;

-- name: GetDefaultOrganisation :one
SELECT *
FROM organisations
WHERE default_org = 1
ORDER BY updated_at DESC, id DESC
LIMIT 1;

-- name: UpdateOrganisation :one
UPDATE organisations
SET friendly_name = ?,
    namespace     = ?,
    default_org   = ?,
    updated_at    = unixepoch('now')
WHERE id = ?
RETURNING *;

-- name: SetOrgsDefaultFalse :exec
UPDATE organisations
SET default_org = false
WHERE default_org = true;

-- name: SetDefaultOrg :one
UPDATE organisations
SET default_org = true
WHERE id = ?
RETURNING *;

-- name: DeleteOrg :exec
DELETE FROM organisations
WHERE id = ?;

-- name: GetOrganisationByID :one
SELECT *
FROM organisations
WHERE id = ?;