-- name: CreateURL :one
INSERT INTO urls(original_url,
                 short_code,
                 is_custom,
                 expired_time)
VALUES ($1, $2, $3, $4)
returning *;

-- name: IsShortCodeAvailable :one
SELECT NOT EXISTS(
    SELECT 1 FROM urls
             WHERE short_code=$1
)AS is_available;

-- name: GetURLByShortCode :one
SELECT * FROM urls
where short_code=$1
AND expired_time>CURRENT_TIMESTAMP;


-- name: DeleteURLExpired :exec
DELETE FROM urls
WHERE expired_time<CURRENT_TIMESTAMP;