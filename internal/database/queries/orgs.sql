-- name: CreateOrganisation :one
INSERT INTO organisations (
  friendly_name,
  namespace,
  created_at,
  updated_at
) VALUES (
  ?,
  ?,
  ?,
  ?
)
RETURNING *;

-- name: ListOrganisations :many
SELECT * FROM organisations
ORDER BY friendly_name;

-- name: GetDefaultOrganisation :one
SELECT * FROM organisations
WHERE default_org = 1
ORDER BY updated_at DESC, id DESC
LIMIT 1;

-- name: UpdateOrganisation :one
UPDATE organisations
SET
  friendly_name = ?,
  namespace = ?,
  default_org = ?,
  updated_at = unixepoch('now')
WHERE id = ?
RETURNING *;