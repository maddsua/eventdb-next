create table streams (
	id text unique primary key,
	push_key text,
	name text not null,
	created_at integer not null default (unixepoch()),
	updated_at integer not null  default (unixepoch())
);

create table events (
	id text unique primary key,
	created_at integer not null default (unixepoch()),
	stream_id text not null,
	client_ip text,
	transaction_id text,
	type text not null,
	level text,
	message text,
	fields text,

	foreign key(stream_id) references streams(id) on update cascade on delete cascade
);
