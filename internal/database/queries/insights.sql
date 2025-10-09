-- name: GetPullRequestInsights :one
WITH average_days_to_merge AS (
    SELECT ROUND((merged_at - created_at) / 86400.0, 2) AS avg_days_to_merge
    FROM pull_requests
    WHERE state = 'MERGED'
      AND created_at >= strftime('%s', 'now', ?)
)
SELECT
    ROUND(MIN(avg_days_to_merge),2) AS min_days_to_merge,
    ROUND(MAX(avg_days_to_merge),2) AS max_days_to_merge,
    ROUND(AVG(avg_days_to_merge),2) AS avg_days_to_merge,
    COUNT(*) AS merged,
    (SELECT COUNT(*) FROM pull_requests WHERE state = 'CLOSED' AND created_at >= strftime('%s', 'now', ?)) AS closed,
    (SELECT COUNT(*) FROM pull_requests WHERE state = 'OPEN' AND created_at >= strftime('%s', 'now', ?)) AS open
FROM average_days_to_merge;

-- name: GetSecuritiesInsights :one
WITH average_days_to_fix AS (
    SELECT ROUND((fixed_at - created_at) / 86400, 2) as days_to_fix
    FROM securities
    WHERE state = 'FIXED'
      AND fixed_at IS NOT NULL
      AND created_at >= strftime('%s', 'now', ?)
)
SELECT
    ROUND(MIN(days_to_fix), 2) AS min_days_to_fix,
    ROUND(MAX(days_to_fix), 2) AS max_days_to_fix,
    ROUND(AVG(days_to_fix), 2) AS avg_days_to_fix,
    (SELECT COUNT(*) FROM securities WHERE state = 'FIXED' AND created_at >= strftime('%s', 'now', ?)) AS fixed,
    (SELECT COUNT(*) FROM securities WHERE state = 'OPEN' AND created_at >= strftime('%s', 'now', ?)) AS open

FROM average_days_to_fix;
