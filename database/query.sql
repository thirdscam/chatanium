-- name: GetUser :one
SELECT id, username, created_at, deleted_at FROM users WHERE id = ?;

-- name: InsertUser :exec
INSERT INTO users (id, username, created_at)
VALUES (?, ?, ?);

-- name: GetGuild :one
SELECT id, name, owner_id FROM guilds WHERE id = ?;

-- name: InsertGuild :exec
INSERT INTO guilds (id, name, owner_id)
VALUES (?, ?, ?)
ON CONFLICT(id) DO UPDATE SET
    name = excluded.name,
    owner_id = excluded.owner_id;

-- name: GetGuildUser :one
SELECT uuid, guild_id, user_id, created_at, quit_at, nickname FROM guildusers
WHERE user_id = ? AND guild_id = ?;

-- name: InsertGuildUser :exec
INSERT INTO guildusers (uuid, guild_id, user_id, created_at, nickname)
VALUES (?, ?, ?, ?, ?);

-- name: GetChannel :one
SELECT id, guild_id, name, description, created_at, deleted_at FROM channels WHERE id = ?;

-- name: InsertChannel :exec
INSERT INTO channels (id, guild_id, name, created_at)
VALUES (?, ?, ?, ?);