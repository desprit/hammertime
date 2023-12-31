// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.23.0
// source: query.sql

package subscription_storage

import (
	"context"
)

const cancelSubscription = `-- name: CancelSubscription :exec
DELETE FROM
    subscription
WHERE
    user_id = ?
    AND schedule_id = ?
`

type CancelSubscriptionParams struct {
	UserID     int64 `json:"user_id"`
	ScheduleID int64 `json:"schedule_id"`
}

func (q *Queries) CancelSubscription(ctx context.Context, arg CancelSubscriptionParams) error {
	_, err := q.db.ExecContext(ctx, cancelSubscription, arg.UserID, arg.ScheduleID)
	return err
}

const createSubscription = `-- name: CreateSubscription :one
INSERT INTO
    subscription (user_id, schedule_id)
VALUES
    (?, ?) RETURNING id, user_id, schedule_id
`

type CreateSubscriptionParams struct {
	UserID     int64 `json:"user_id"`
	ScheduleID int64 `json:"schedule_id"`
}

func (q *Queries) CreateSubscription(ctx context.Context, arg CreateSubscriptionParams) (Subscription, error) {
	row := q.db.QueryRowContext(ctx, createSubscription, arg.UserID, arg.ScheduleID)
	var i Subscription
	err := row.Scan(&i.ID, &i.UserID, &i.ScheduleID)
	return i, err
}

const deleteSubscription = `-- name: DeleteSubscription :exec
DELETE FROM
    subscription
WHERE
    id = ?
`

func (q *Queries) DeleteSubscription(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteSubscription, id)
	return err
}

const getSubscription = `-- name: GetSubscription :one
SELECT
    id, user_id, schedule_id
FROM
    subscription
WHERE
    user_id = ?
    AND schedule_id = ?
LIMIT
    1
`

type GetSubscriptionParams struct {
	UserID     int64 `json:"user_id"`
	ScheduleID int64 `json:"schedule_id"`
}

func (q *Queries) GetSubscription(ctx context.Context, arg GetSubscriptionParams) (Subscription, error) {
	row := q.db.QueryRowContext(ctx, getSubscription, arg.UserID, arg.ScheduleID)
	var i Subscription
	err := row.Scan(&i.ID, &i.UserID, &i.ScheduleID)
	return i, err
}

const getSubscriptionById = `-- name: GetSubscriptionById :one
SELECT
    id, user_id, schedule_id
FROM
    subscription
WHERE
    id = ?
LIMIT
    1
`

func (q *Queries) GetSubscriptionById(ctx context.Context, id int64) (Subscription, error) {
	row := q.db.QueryRowContext(ctx, getSubscriptionById, id)
	var i Subscription
	err := row.Scan(&i.ID, &i.UserID, &i.ScheduleID)
	return i, err
}

const getSubscriptions = `-- name: GetSubscriptions :many
SELECT
    id, user_id, schedule_id
FROM
    subscription
`

func (q *Queries) GetSubscriptions(ctx context.Context) ([]Subscription, error) {
	rows, err := q.db.QueryContext(ctx, getSubscriptions)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Subscription
	for rows.Next() {
		var i Subscription
		if err := rows.Scan(&i.ID, &i.UserID, &i.ScheduleID); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
