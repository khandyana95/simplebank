-- name: CreateEntry :one
INSERT INTO entries (
  account_id, amount, currency
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: GetEntry :one
SELECT * FROM entries
WHERE id = $1 LIMIT 1;

-- name: ListEntriesByAccountId :many
SELECT * FROM entries
WHERE account_id = $1
ORDER by id
LIMIT $2
OFFSET $3;

-- name: DeleteEntry :exec
DELETE FROM entries
WHERE id = $1;