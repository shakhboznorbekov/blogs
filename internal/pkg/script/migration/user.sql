create table users
(
    id         text primary key not null ,
    username   text not null ,
    password   text,
    role       text not null ,
    status     boolean not null ,
    created_at timestamp default now(),
    deleted_at timestamp,
    updated_at timestamp,
    updated_by text references users(id),
    created_by text references users(id),
    deleted_by text references users(id)
);

alter table users
    owner to postgres;

