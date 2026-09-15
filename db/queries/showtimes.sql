-- name: CreateShowtime :one
INSERT INTO showtimes (
    id, movie_id, studio_id, start_time, end_time, price, status
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
)
RETURNING id, movie_id, studio_id, start_time, end_time, price, status, created_at, updated_at;

-- name: GetShowtimeByID :one
SELECT 
    s.id, s.movie_id, s.studio_id, s.start_time, s.end_time, s.price, s.status, s.created_at, s.updated_at,
    m.title AS movie_title, m.duration_minutes AS movie_duration, m.rating AS movie_rating,
    st.name AS studio_name,
    c.name AS cinema_name
FROM showtimes s
JOIN movies m ON s.movie_id = m.id
JOIN studios st ON s.studio_id = st.id
JOIN cinemas c ON st.cinema_id = c.id
WHERE s.id = $1 LIMIT 1;

-- name: ListShowtimes :many
SELECT 
    s.id, s.movie_id, s.studio_id, s.start_time, s.end_time, s.price, s.status, s.created_at, s.updated_at,
    m.title AS movie_title, m.duration_minutes AS movie_duration, m.rating AS movie_rating,
    st.name AS studio_name,
    c.name AS cinema_name
FROM showtimes s
JOIN movies m ON s.movie_id = m.id
JOIN studios st ON s.studio_id = st.id
JOIN cinemas c ON st.cinema_id = c.id
WHERE s.status != 'CANCELLED'
ORDER BY s.start_time ASC
LIMIT $1 OFFSET $2;

-- name: UpdateShowtime :one
UPDATE showtimes
SET 
    movie_id = $2,
    studio_id = $3,
    start_time = $4,
    end_time = $5,
    price = $6,
    status = $7,
    updated_at = NOW()
WHERE id = $1
RETURNING id, movie_id, studio_id, start_time, end_time, price, status, created_at, updated_at;

-- name: DeleteShowtime :exec
UPDATE showtimes
SET 
    status = 'CANCELLED',
    updated_at = NOW()
WHERE id = $1;

-- name: CheckShowtimeOverlap :one
SELECT COUNT(*) FROM showtimes
WHERE studio_id = $1
  AND id != $2
  AND status != 'CANCELLED'
  AND (start_time < $4 AND end_time > $3);