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

-- name: GetFeedFollowsForUser :many
select
    feed_follows.feed_id,
    feed_follows.user_id,
    feeds.name as feed_name,
    feeds.url as feed_url
from feed_follows
    inner join feeds on feeds.id = feed_follows.feed_id where feed_follows.user_id = $1;

-- name: DeleteFollow :exec
delete from feed_follows where user_id = $1 and feed_id = $2;
