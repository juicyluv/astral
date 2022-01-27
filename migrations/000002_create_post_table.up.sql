CREATE TABLE IF NOT EXISTS posts(
    post_id serial primary key not null,
    title text not null,
    content text not null,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    author_id int not null,

    foreign key(author_id) references users(user_id)
);