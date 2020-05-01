
-- list of all current metrics names
-- name: GetMetrics :many
select distinct(m.source)
from public.metrics m;


-- name: MetricForSource :many
select *
from public.metrics
where source = $1;

-- name: MetricStatsPerPeriod :many
select m.source,
       m.name,
       to_timestamp(floor((extract('epoch' from m.ts) / sqlc.arg(seconds)::int)) * sqlc.arg(seconds)::int) ts,
       avg(m.value) avg,
       max(m.value) max,
       min(m.value) min
from public.metrics m
group by m.source, m.name, ts_bucket;

