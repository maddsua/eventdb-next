-- name: CreateStream :execrows
insert into streams (
	id,
	push_key,
	name
) values (
	sqlc.arg(id),
	sqlc.narg(push_key),
	sqlc.arg(name)
);

-- name: GetStreamByID :one
select * from streams
where id = sqlc.arg(id);

-- name: GetStreams :many
select * from streams
order by created_at asc;

-- name: SetStreamPushKey :execrows
update streams
	set push_key = sqlc.narg(push_key)
where id = sqlc.arg(id);

-- name: SetStreamName :execrows
update streams
	set name = sqlc.arg(name)
where id = sqlc.arg(id);

-- name: RemoveStream :execrows
delete from streams
where id = sqlc.arg(id);

-- name: AddEvent :exec
insert into events (
	id,
	created_at,
	stream_id,
	client_ip,
	transaction_id,
	type,
	level,
	message,
	fields
) values (
	sqlc.arg(id),
	sqlc.narg(created_at),
	sqlc.arg(stream_id),
	sqlc.narg(client_ip),
	sqlc.narg(transaction_id),
	sqlc.narg(type),
	sqlc.narg(level),
	sqlc.narg(message),
	sqlc.narg(fields)
);

-- name: DeleteEvents :many
delete from events
where (id = sqlc.narg(id) or sqlc.narg(id) is null)
	and (stream_id = sqlc.narg(stream_id) or sqlc.narg(stream_id) is null)
	and (created_at < sqlc.narg(before) or sqlc.narg(before) is null)
	and (created_at > sqlc.narg(after) or sqlc.narg(after) is null)
returning id;

-- name: GetEvents :many
select * from events
where (stream_id = sqlc.narg(stream_id) or sqlc.narg(stream_id) is null)
	and (level = sqlc.narg(level) or sqlc.narg(level) is null)
	and (created_at < sqlc.narg(before) or sqlc.narg(before) is null)
	and (created_at > sqlc.narg(after) or sqlc.narg(after) is null);

-- name: GetActivity :many
select
	created_at,
	level
from events
where (created_at < sqlc.narg(before) or sqlc.narg(before) is null)
	and (created_at > sqlc.narg(after) or sqlc.narg(after) is null);
