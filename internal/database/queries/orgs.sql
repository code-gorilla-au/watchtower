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
