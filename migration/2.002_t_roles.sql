CREATE TABLE public.roles (
	id int8 DEFAULT int64_id('roles'::text, 'id'::text) NOT NULL,
	"name" varchar(45) NOT NULL,
	CONSTRAINT roles_pk PRIMARY KEY (id),
	CONSTRAINT roles_unique UNIQUE (name)
);