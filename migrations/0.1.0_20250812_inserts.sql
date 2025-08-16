insert into permissions(name) values ('PERM');
insert into roles(name) values ('ROLE_1');
insert into roles(name) values ('ROLE_2');
insert into permissions_roles (permission_id, role_id) values (1, 1);
insert into permissions_roles (permission_id, role_id) values (1, 2);

insert into users (code, username, email, password, salt, status, created_by, create_date, permission_id) values ('e60d5b12-921d-4f1e-966a-2d4e743a164b', 'javi', 'javi@mail.com', 'vIogDylSzy647bVGSNea3RVie+lcb0Ex1b6nTogYOck=', '6/xakjp96vErhv+PD6xVqA==', 'ACTIVE', 'javi', now(), 1);

insert into users (code, username, email, password, salt, status, created_by, create_date, permission_id) values ('e60d5b12-921d-4f1e-966a-2d4e743a164a', 'johan', 'johan@mail.com', 'vIogDylSzy647bVGSNea3RVie+lcb0Ex1b6nTogYOck=', '6/xakjp96vErhv+PD6xVqA==', 'ACTIVE', 'javi', now(), 1);
