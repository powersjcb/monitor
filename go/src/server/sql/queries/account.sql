-- name: InsertAccount :one
INSERT INTO public.accounts (auth_provider, auth_provider_id, api_key, inserted_at)
VALUES ($1, $2, $3, NOW())
RETURNING *;

-- name: GetAccountByID :one
SELECT *
FROM public.accounts
WHERE id = $1
LIMIT 1;

-- name: GetAccountByProviderID :one
SELECT *
FROM public.accounts
WHERE auth_provider = $1 AND auth_provider_id = $2
LIMIT 1;


-- name: GetAccountIDForAPIKey :one
SELECT id
FROM public.accounts
WHERE api_key = $1;
