CREATE TABLE IF NOT EXISTS users(
    id serial primary key not null,
    username text not null,
    email text not null,
    registered_at timestamptz not null default now(),
    password text
);