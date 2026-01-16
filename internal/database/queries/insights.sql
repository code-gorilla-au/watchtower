-- name: GetPullRequestInsightsByOrg :one
WITH
    pull_request_with_org AS (
        SELECT pr.*, po.organisation_id as organisation_id, p.id as product_id
        FROM pull_requests pr
                 JOIN repositories r ON r.name = pr.repository_name
                 JOIN product_organisations po
                 JOIN products p ON p.id = po.product_id
            AND JSON_VALID(p.tags)
            AND EXISTS (SELECT 1
                        FROM JSON_EACH(p.tags)
                        WHERE JSON_EACH.value = r.topic)
        WHERE po.organisation_id = ?
        ORDER BY pr.created_at DESC
    ),
    average_days_to_merge AS (
        SELECT ROUND((merged_at - created_at) / 86400.0, 2) AS avg_days_to_merge
        FROM pull_request_with_org
        WHERE state = 'MERGED'
          AND created_at >= strftime('%s', 'now', ?)
    )
SELECT
    ROUND(COALESCE(MIN(avg_days_to_merge), 0), 2) AS min_days_to_merge,
    ROUND(COALESCE(MAX(avg_days_to_merge), 0), 2) AS max_days_to_merge,
    ROUND(COALESCE(AVG(avg_days_to_merge), 0), 2) AS avg_days_to_merge,
    COUNT(avg_days_to_merge) AS merged,
    (SELECT COUNT(*) FROM pull_request_with_org WHERE state = 'CLOSED' AND created_at >= strftime('%s', 'now', ?)) AS closed,
    (SELECT COUNT(*) FROM pull_request_with_org WHERE state = 'OPEN' AND created_at >= strftime('%s', 'now', ?)) AS open
FROM average_days_to_merge;

-- name: GetSecuritiesInsightsByOrg :one
WITH
    securities_with_org AS (
        SELECT s.*, po.organisation_id as organisation_id, p.id as product_id
        FROM securities s
                 JOIN repositories r ON r.name = s.repository_name
                 JOIN product_organisations po
                 JOIN products p ON p.id = po.product_id
            AND JSON_VALID(p.tags)
            AND EXISTS (SELECT 1
                        FROM JSON_EACH(p.tags)
                        WHERE JSON_EACH.value = r.topic)
        WHERE po.organisation_id = ?
        ORDER BY s.created_at DESC
    ),
    average_days_to_fix AS (
        SELECT ROUND((fixed_at - created_at) / 86400.0, 2) as days_to_fix
        FROM securities_with_org
        WHERE state = 'FIXED'
          AND fixed_at IS NOT NULL
          AND created_at >= strftime('%s', 'now', ?)
    )
SELECT
    ROUND(COALESCE(MIN(days_to_fix), 0), 2) AS min_days_to_fix,
    ROUND(COALESCE(MAX(days_to_fix), 0), 2) AS max_days_to_fix,
    ROUND(COALESCE(AVG(days_to_fix), 0), 2) AS avg_days_to_fix,
    (SELECT COUNT(*) FROM securities_with_org WHERE state = 'FIXED' AND created_at >= strftime('%s', 'now', ?)) AS fixed,
    (SELECT COUNT(*) FROM securities_with_org WHERE state = 'OPEN' AND created_at >= strftime('%s', 'now', ?)) AS open
FROM average_days_to_fix;