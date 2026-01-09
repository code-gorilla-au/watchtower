-- name: CreateOrgNotification :one
INSERT INTO organisation_notifications (organisation_id,
                                        external_id,
                                        type,
                                        content,
                                        created_at,
                                        updated_at)
VALUES (?,
        ?,
        ?,
        ?,
        CAST(strftime('%s', 'now') AS INTEGER),
        CAST(strftime('%s', 'now') AS INTEGER))
RETURNING *;

-- name: GetNotificationByExternalID :one
SELECT *
FROM organisation_notifications
WHERE external_id = ?;

-- name: UpdateOrgNotificationByID :one
UPDATE organisation_notifications
SET type       = ?,
    content    = ?,
    status     = ?,
    updated_at = CAST(strftime('%s', 'now') AS INTEGER)
WHERE id = ?
RETURNING *;

-- name: UpdateOrgNotificationStatusByID :one
UPDATE organisation_notifications
SET status     = ?,
    updated_at = CAST(strftime('%s', 'now') AS INTEGER)
WHERE id = ?
RETURNING *;

-- name: GetUnreadNotifications :many
SELECT *
FROM organisation_notifications
WHERE status = 'unread'
ORDER BY created_at DESC;

-- name: DeleteOrgNotificationByDate :exec
DELETE
FROM organisation_notifications
WHERE created_at < ?;
