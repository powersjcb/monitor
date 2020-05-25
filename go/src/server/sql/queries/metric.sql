-- name: InsertMetric :one
INSERT INTO public.metrics (ts, source, name, target, value, inserted_at)
VALUES ($1, $2, $3, $4, $5, NOW())
RETURNING *;

-- name: GetMetricStatsPerPeriod :many
select m.source,
       m.name,
       to_timestamp(floor((extract('epoch' from m.ts) / sqlc.arg(seconds)::int)) * sqlc.arg(seconds)::int) ts,
       avg(m.value) avg,
       max(m.value) max,
       min(m.value) min
from public.metrics m
where account_id = sqlc.arg(account_id)::bigint
group by m.source, m.name, ts_bucket;

