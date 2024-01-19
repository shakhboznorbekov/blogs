create table loggers(
                        id serial primary key,
                        created_at timestamp default now(),
                        data jsonb,
                        method text,
                        action text
);