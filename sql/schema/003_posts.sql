-- +goose Up
CREATE TABLE posts (
    id uuid unique primary key,
    created_at timestamp not null,
    updated_at timestamp not null,
    title text not null,
    url text unique not null,
    description text unique not null,
    published_at timestamp not null,
    feed_id uuid not null,
    foreign key (feed_id) references feeds(id)
);

-- +goose Down
DROP TABLE posts;
