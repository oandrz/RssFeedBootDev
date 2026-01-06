-- +goose Up
CREATE TABLE feeds (
   name text not null,
   url text unique not null,
   user_id uuid not null,
   foreign key (user_id) references users(id) on delete cascade
);

-- +goose Down
DROP TABLE feeds;
