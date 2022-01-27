CREATE TABLE IF NOT EXISTS user_post(
    user_id int not null,
    post_id int not null,

    foreign key(user_id) references users(user_id) on delete cascade,
    foreign key(post_id) references posts(post_id) on delete cascade
);