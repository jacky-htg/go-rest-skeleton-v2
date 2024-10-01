CREATE TABLE public.access_roles (
	access_id int8 NOT NULL,
	role_id int8 NOT NULL,
	CONSTRAINT access_roles_pk PRIMARY KEY (access_id, role_id)
);