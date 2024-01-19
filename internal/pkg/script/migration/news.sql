CREATE TABLE  news(
    id            text primary key not null ,
    title VARCHAR(255) NOT NULL,
    content TEXT,
    source VARCHAR(100),
    created_at    timestamp default now(),
    deleted_at    timestamp,
    updated_at    timestamp,
    updated_by    text references users(id),
    created_by    text references users(id),
    deleted_by    text references users(id)
    );