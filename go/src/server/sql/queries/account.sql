
-- name: GetOrCreateAccount :one
INSERT INTO public.accounts (auth_provider_id, auth_provider, api_key, inserted_at)
VALUES ($1, $2, $3, NOW())
ON CONFLICT (auth_provider_id, auth_provider) DO NOTHING
RETURNING *;

-- name: GetAccountIDForAPIKey :one
SELECT id
FROM public.accounts
WHERE api_key = $1;
