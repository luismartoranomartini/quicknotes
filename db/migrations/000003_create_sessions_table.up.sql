create table sessions (
	token TEXT PRIMARY KEY,
	data BYTEA not null,
	expiry TIMESTAMPTZ NOT NULL
);

create index sessions_expiry_idx on sessions (expiry);
