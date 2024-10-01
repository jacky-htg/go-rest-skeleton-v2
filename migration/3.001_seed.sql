INSERT INTO public.users (id,"name",email,"password",created_at,created_by) VALUES
	 (425071490427828,'Rijal Asepnugroho','rijal.asep.nugroho@gmail.com','$2a$10$eCZXQBWquZJrlKglS8trh.5l2UnM8.m0Sah4T0iHE5QBgeJov2kBO','2024-09-26 14:18:28.946672+07',425071490427828);

INSERT INTO public.roles (id,"name") VALUES(156677038157782,'Superman');

INSERT INTO public."access" (id,"name","path") VALUES
	 (495991231925511,'list user','GET /users'),
	 (144157733335917,'create user','POST /users'),
	 (852228553691053,'view user','GET /users/:id'),
	 (697344866789774,'update user','PUT /users/:id'),
	 (119395353616382,'delete user','DELETE /users/:id');

INSERT INTO public.roles_users (user_id,role_id) VALUES (425071490427828,156677038157782);

INSERT INTO public.access_roles (access_id,role_id) VALUES
	 (495991231925511,156677038157782),
	 (144157733335917,156677038157782),
	 (852228553691053,156677038157782),
	 (697344866789774,156677038157782),
	 (119395353616382,156677038157782);

