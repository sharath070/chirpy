-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (
    token,
    created_at,
    updated_at,
    expires_at,
    revoked_at,
    user_id
) VALUES (
    $1,
    NOW(),
    NOW(),
    NOW() + INTERVAL '60 days',
    NULL,
    $2
)
RETURNING *;

-- name: GetUserFromRefreshToken :one
SELECT 
    users.id AS user_id,  
    users.email AS user_email
FROM users
INNER JOIN refresh_tokens on users.id = refresh_tokens.user_id
WHERE
    refresh_tokens.token = $1
    AND (refresh_tokens.expires_at > NOW())
    AND (refresh_tokens.revoked_at IS NULL);
