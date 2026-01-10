-- +goose Up
CREATE TABLE feeds (
   id uuid Primary Key,
   name text not null,
   url text unique not null,
   last_fetched_at timestamp,
   user_id uuid not null,
   foreign key (user_id) references users(id) on delete cascade
);

create table feed_follows (
    id uuid primary key,
    created_at timestamp not null,
    updated_at timestamp not null,
    feed_id uuid not null,
    user_id uuid not null,
    foreign key (feed_id) references feeds(id) on delete cascade,
    foreign key (user_id) references users(id) on delete cascade,
    unique (user_id, feed_id)
);

-- +goose Down
DROP TABLE feed_follows;
DROP TABLE feeds;
