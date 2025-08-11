drop table if exists users;
drop table if exists permissions_roles;
drop table if exists permissions;
drop table if exists roles;

create table permissions (
    id serial primary key,
    name text unique not null
);

create table roles (
    id serial primary key,
    name text unique not null
);

create table permissions_roles (
    permission_id integer not null,
    role_id integer not null,
    primary key(permission_id, role_id),
    foreign key(permission_id) references permissions(id),
    foreign key(role_id) references roles(id)
);

create table users (
    id serial primary key,
    code uuid not null,
    username text not null,
    email text not null,
    password text not null,
    salt text not null,
    status text not null,
    created_by text not null,
    last_modified_by text,
    create_date date not null,
    last_modified_date date,
    permission_id integer not null,
    foreign key(permission_id) references permissions(id)
);
