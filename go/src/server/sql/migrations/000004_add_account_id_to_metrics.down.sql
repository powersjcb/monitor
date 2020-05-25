-- down

DROP INDEX metrics_account_id_ts;
create index source_ts_index on metrics (source, ts);
