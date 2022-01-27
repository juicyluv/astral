CREATE TABLE IF NOT EXISTS posts(
    id serial primary key not null,
    title text not null,
    subtitle text not null,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    author_id int not null,

    foreign key(author_id) references users(user_id)
);