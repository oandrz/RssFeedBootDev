-- name: CreateFeed :one
INSERT INTO feeds (id, name, url, user_id)
VALUES (
        $1,
   $2,
   $3,
   $4
)
RETURNING *;

-- name: GetFeeds :many
select * from feeds;

-- name: CreateFeedFollow :one
with inserted_feed_follow as (
    INSERT INTO feed_follows (id, created_at, updated_at, feed_id, user_id)
    VALUES (
        $1,
        $2,
        $3,
        $4,
        $5
    )
    RETURNING *
)
select
    inserted_feed_follow.*,
    feeds.name as feed_name,
    users.name as user_name
from inserted_feed_follow
         inner join feeds on feeds.id = inserted_feed_follow.feed_id
         inner join users on users.id = inserted_feed_follow.user_id;

-- name: GetFeedByUrl :one
select * from feeds where url = $1;