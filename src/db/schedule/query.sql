-- name: GetScheduleEntry :one
SELECT
    *
FROM
    schedule
WHERE
    id = ?
LIMIT
    1;

-- name: GetScheduleEntryByActivityID :one
SELECT
    *
FROM
    schedule
WHERE
    activity_id = ?
LIMIT
    1;

-- name: GetScheduleEntriesWithPreEntry :many
SELECT
    *
FROM
    schedule
WHERE
    datetime > datetime('now')
    AND begin_date > datetime('now')
    AND pre_entry = 1
ORDER BY
    datetime ASC;

-- name: GetLatestScheduleEntry :one
SELECT
    *
FROM
    schedule
ORDER BY
    datetime DESC
LIMIT
    1;

-- name: CreateScheduleEntry :one
INSERT
    or REPLACE INTO schedule (
        activity_id,
        datetime,
        trainer,
        activity,
        pre_entry,
        begin_date
    )
VALUES
    (?, ?, ?, ?, ?, ?) RETURNING *;

-- name: UpdateScheduleEntry :exec
UPDATE
    schedule
SET
    activity_id = ?,
    datetime = ?,
    trainer = ?,
    activity = ?,
    pre_entry = ?,
    begin_date = ?
WHERE
    ID = ?;

-- name: DeleteScheduleEntry :exec
DELETE FROM
    schedule
WHERE
    ID = ?;