CREATE TABLE users (
	id int8 DEFAULT int64_id('users'::text, 'id'::text) NOT NULL,
	"name" varchar(45) NOT NULL,
	email varchar(128) NOT NULL,
	"password" bpchar(60) NOT NULL,
	created_at timestamptz DEFAULT timezone('utc'::text, now()) NULL,
	created_by int8 NOT NULL,
	updated_at timestamptz NULL,
	updated_by int8 NULL,
	deleted_at timestamptz NULL,
	deleted_by int8 NULL,
	CONSTRAINT users_email_key UNIQUE (email),
	CONSTRAINT users_pkey PRIMARY KEY (id)
);