-- down

ALTER TABLE public.metrics DROP CONSTRAINT fk_metric_account;

ALTER TABLE public.metrics ALTER COLUMN account_id DROP NOT NULL;