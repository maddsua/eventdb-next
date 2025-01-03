// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: queries.sql

package sqliteops

import (
	"context"
	"database/sql"
)

const addEvent = `-- name: AddEvent :exec
insert into events (
	id,
	stream_id,
	client_ip,
	transaction_id,
	type,
	level,
	message,
	fields
) values (
	?1,
	?2,
	?3,
	?4,
	?5,
	?6,
	?7,
	?8
)
`

type AddEventParams struct {
	ID            string
	StreamID      string
	ClientIp      sql.NullString
	TransactionID sql.NullString
	Type          sql.NullString
	Level         sql.NullString
	Message       sql.NullString
	Fields        sql.NullString
}

func (q *Queries) AddEvent(ctx context.Context, arg AddEventParams) error {
	_, err := q.db.ExecContext(ctx, addEvent,
		arg.ID,
		arg.StreamID,
		arg.ClientIp,
		arg.TransactionID,
		arg.Type,
		arg.Level,
		arg.Message,
		arg.Fields,
	)
	return err
}

const createStream = `-- name: CreateStream :execrows
insert into streams (
	id,
	push_key,
	name
) values (
	?1,
	?2,
	?3
)
`

type CreateStreamParams struct {
	ID      string
	PushKey sql.NullString
	Name    string
}

func (q *Queries) CreateStream(ctx context.Context, arg CreateStreamParams) (int64, error) {
	result, err := q.db.ExecContext(ctx, createStream, arg.ID, arg.PushKey, arg.Name)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

const deleteEvents = `-- name: DeleteEvents :execrows
delete from events
where (id = ?1 or ?1 is null)
	and (stream_id = ?2 or ?2 is null)
	and (created_at < ?3 or ?3 is null)
	and (created_at > ?4 or ?4 is null)
`

type DeleteEventsParams struct {
	ID       sql.NullString
	StreamID sql.NullString
	Before   sql.NullInt64
	After    sql.NullInt64
}

func (q *Queries) DeleteEvents(ctx context.Context, arg DeleteEventsParams) (int64, error) {
	result, err := q.db.ExecContext(ctx, deleteEvents,
		arg.ID,
		arg.StreamID,
		arg.Before,
		arg.After,
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

const getEvents = `-- name: GetEvents :many
select id, created_at, stream_id, client_ip, transaction_id, type, level, message, fields from events
where (stream_id = ?1 or ?1 is null)
	and (level = ?2 or ?2 is null)
	and (created_at < ?3 or ?3 is null)
	and (created_at > ?4 or ?4 is null)
`

type GetEventsParams struct {
	StreamID sql.NullString
	Level    sql.NullString
	Before   sql.NullInt64
	After    sql.NullInt64
}

func (q *Queries) GetEvents(ctx context.Context, arg GetEventsParams) ([]Event, error) {
	rows, err := q.db.QueryContext(ctx, getEvents,
		arg.StreamID,
		arg.Level,
		arg.Before,
		arg.After,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Event
	for rows.Next() {
		var i Event
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.StreamID,
			&i.ClientIp,
			&i.TransactionID,
			&i.Type,
			&i.Level,
			&i.Message,
			&i.Fields,
		); err != nil {
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

const getStreamByID = `-- name: GetStreamByID :one
select id, push_key, name, created_at, updated_at from streams
where id = ?1
`

func (q *Queries) GetStreamByID(ctx context.Context, id string) (Stream, error) {
	row := q.db.QueryRowContext(ctx, getStreamByID, id)
	var i Stream
	err := row.Scan(
		&i.ID,
		&i.PushKey,
		&i.Name,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getStreams = `-- name: GetStreams :many
select id, push_key, name, created_at, updated_at from streams
order by created_at asc
limit ?2 offset ?1
`

type GetStreamsParams struct {
	Offset int64
	Limit  int64
}

func (q *Queries) GetStreams(ctx context.Context, arg GetStreamsParams) ([]Stream, error) {
	rows, err := q.db.QueryContext(ctx, getStreams, arg.Offset, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Stream
	for rows.Next() {
		var i Stream
		if err := rows.Scan(
			&i.ID,
			&i.PushKey,
			&i.Name,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
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

const removeStream = `-- name: RemoveStream :execrows
delete from streams
where id = ?1
`

func (q *Queries) RemoveStream(ctx context.Context, id string) (int64, error) {
	result, err := q.db.ExecContext(ctx, removeStream, id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

const setStreamName = `-- name: SetStreamName :execrows
update streams
	set name = ?1
where id = ?2
`

type SetStreamNameParams struct {
	Name string
	ID   string
}

func (q *Queries) SetStreamName(ctx context.Context, arg SetStreamNameParams) (int64, error) {
	result, err := q.db.ExecContext(ctx, setStreamName, arg.Name, arg.ID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

const setStreamPushKey = `-- name: SetStreamPushKey :execrows
update streams
	set push_key = ?1
where id = ?2
`

type SetStreamPushKeyParams struct {
	PushKey sql.NullString
	ID      string
}

func (q *Queries) SetStreamPushKey(ctx context.Context, arg SetStreamPushKeyParams) (int64, error) {
	result, err := q.db.ExecContext(ctx, setStreamPushKey, arg.PushKey, arg.ID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
