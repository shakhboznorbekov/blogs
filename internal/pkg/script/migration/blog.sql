CREATE TABLE  blogs(
    id            text primary key not null ,
    title         varchar(255) NOT NULL,
    content       text,
    author        varchar(100),
    created_at    timestamp default now(),
    deleted_at    timestamp,
    updated_at    timestamp,
    updated_by    text references users(id),
    created_by    text references users(id),
    deleted_by    text references users(id)
    );