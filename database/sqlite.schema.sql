create table streams (
	id text unique primary key,
	push_key text not null,
	name text not null,
	created_at integer not null default (unixepoch()),
	updated_at integer not null  default (unixepoch())
);

create table events (
	id text unique primary key,
	date integer not null default (unixepoch()),
	stream_id text not null,
	client_ip text null,
	request_id text null,
	type text not null,
	level text null,
	message text null,
	fields text null,

	foreign key(stream_id) references streams(id) on update cascade on delete cascade
);
