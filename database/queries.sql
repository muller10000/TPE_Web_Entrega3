-- name: CreateMovie :one
INSERT INTO movies (title, director, year, genre, rating)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetMovie :one
SELECT * FROM movies WHERE id = $1;

-- name: ListMovies :many
SELECT * FROM movies ORDER BY created_at DESC;

-- name: UpdateMovie :one
UPDATE movies
SET title = $1,
    director = $2,
    year = $3,
    genre = $4,
    rating = $5
WHERE id = $6
RETURNING *;

-- name: DeleteMovie :exec
DELETE FROM movies WHERE id = $1;

-- NUEVAS QUERIES: Usuarios
-- name: GetUser :one
SELECT * FROM users WHERE username = $1;

-- name: CreateUser :one
INSERT INTO users (username, password) VALUES ($1, $2) RETURNING *;