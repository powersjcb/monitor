-- up

DROP INDEX source_ts_index;

CREATE INDEX metrics_account_id_ts ON public.metrics (account_id, ts);
