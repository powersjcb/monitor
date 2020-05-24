-- up

create table public.accounts (
    id bigserial primary key,
    auth_provider_id varchar not null,
    auth_provider varchar not null,
    api_key varchar not null,
    inserted_at timestamp
);

create unique index on public.accounts (auth_provider, auth_provider_id);
create unique index on public.accounts (api_key);

ALTER TABLE public.metrics ADD account_id bigint;