CREATE TABLE public."access" (
	id int8 DEFAULT int64_id('access'::text, 'id'::text) NOT NULL,
	"name" varchar(128) NOT NULL,
	"path" varchar(128) NOT NULL,
	CONSTRAINT newtable_pk PRIMARY KEY (id),
	CONSTRAINT newtable_unique UNIQUE (name),
	CONSTRAINT newtable_unique_1 UNIQUE (path)
);