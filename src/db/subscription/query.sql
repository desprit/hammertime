-- name: GetSubscriptionById :one
SELECT
    *
FROM
    subscription
WHERE
    id = ?
LIMIT
    1;

-- name: GetSubscription :one
SELECT
    *
FROM
    subscription
WHERE
    user_id = ?
    AND schedule_id = ?
LIMIT
    1;

-- name: GetSubscriptions :many
SELECT
    *
FROM
    subscription;

-- name: CreateSubscription :one
INSERT INTO
    subscription (user_id, schedule_id)
VALUES
    (?, ?) RETURNING *;

-- name: CancelSubscription :exec
DELETE FROM
    subscription
WHERE
    user_id = ?
    AND schedule_id = ?;

-- name: DeleteSubscription :exec
DELETE FROM
    subscription
WHERE
    id = ?;