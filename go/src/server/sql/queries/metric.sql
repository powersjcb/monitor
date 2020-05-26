-- name: InsertMetric :one
INSERT INTO public.metrics (account_id, ts, source, name, target, value, inserted_at)
VALUES ($1, $2, $3, $4, $5, $6, NOW())
RETURNING *;

-- name: GetMetricStatsPerPeriod :many
select m.source,
       m.name,
       to_timestamp(floor((extract('epoch' from m.ts) / sqlc.arg(seconds)::int)) * sqlc.arg(seconds)::int)::timestamp ts,
       avg(m.value)::float avg,
       max(m.value)::float max,
       min(m.value)::float min
from public.metrics m
where account_id = sqlc.arg(account_id)::bigint
group by m.source, m.name, ts_bucket;