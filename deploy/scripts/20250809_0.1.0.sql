create table permissions (
    id serial primary key,
    name text not null
);

create table roles (
    id serial primary key,
    name text not null
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
    username text not null,
    email text not null,
    password text not null,
    salt text not null,
    status text not null,
    permission_id integer not null,
    foreign key(permission_id) references permissions(id)
);
