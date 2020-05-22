
-- name: InsertAccount :one
INSERT INTO public.accounts (username, api_key, inserted_at)
    VALUES ($1, $2, NOW())
    RETURNING *;

-- name: GetAccountIDForAPIKey :one
SELECT id
FROM public.accounts
WHERE api_key = $1;
