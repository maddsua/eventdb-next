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

-- name: AddEvent :execrows
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
	sqlc.arg(id),
	sqlc.arg(stream_id),
	sqlc.narg(client_ip),
	sqlc.narg(transaction_id),
	sqlc.narg(type),
	sqlc.narg(level),
	sqlc.narg(message),
	sqlc.narg(fields)
);

-- name: DeleteEvents :execrows
delete from events
where (id = sqlc.narg(id) or sqlc.narg(id) is null)
	and (stream_id = sqlc.narg(stream_id) or sqlc.narg(stream_id) is null)
	and (created_at < sqlc.narg(before) or sqlc.narg(before) is null)
	and (created_at > sqlc.narg(after) or sqlc.narg(after) is null);
