-- up


ALTER TABLE public.metrics
ADD CONSTRAINT fk_metric_account FOREIGN KEY (account_id) REFERENCES accounts(id);

ALTER TABLE public.metrics ALTER COLUMN account_id SET NOT NULL;
