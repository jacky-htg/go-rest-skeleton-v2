CREATE TABLE public.roles_users (
	user_id int8 NOT NULL,
	role_id int8 NOT NULL,
	CONSTRAINT roles_users_pk PRIMARY KEY (user_id, role_id)
);