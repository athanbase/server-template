-- name: QueryUsers :many
SELECT * FROM user;

-- name: CreateUser :execlastid
INSERT INTO user (name) VALUES(?);

-- name: UpdateUser :execrows
UPDATE user SET country = ? WHERE name = ?;