-- up

create table public.accounts (
    id bigserial primary key,
    username varchar not null,
    api_key varchar not null,
    inserted_at timestamp
);

ALTER TABLE public.metrics ADD account_id bigint;