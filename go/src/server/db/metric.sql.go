// Code generated by sqlc. DO NOT EDIT.
// source: metric.sql

package db

import (
	"context"
	"database/sql"
)

const getMetricStatsPerPeriod = `-- name: GetMetricStatsPerPeriod :many
select m.source,
       m.name,
       to_timestamp(floor((extract('epoch' from m.ts) / $1::int)) * $1::int) ts,
       avg(m.value) avg,
       max(m.value) max,
       min(m.value) min
from public.metrics m
where account_id = $2::bigint
group by m.source, m.name, ts_bucket
`

type GetMetricStatsPerPeriodParams struct {
	Seconds   int32 `json:"seconds"`
	AccountID int64 `json:"account_id"`
}

type GetMetricStatsPerPeriodRow struct {
	Source string      `json:"source"`
	Name   string      `json:"name"`
	Ts     interface{} `json:"ts"`
	Avg    interface{} `json:"avg"`
	Max    interface{} `json:"max"`
	Min    interface{} `json:"min"`
}

func (q *Queries) GetMetricStatsPerPeriod(ctx context.Context, arg GetMetricStatsPerPeriodParams) ([]GetMetricStatsPerPeriodRow, error) {
	rows, err := q.db.QueryContext(ctx, getMetricStatsPerPeriod, arg.Seconds, arg.AccountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetMetricStatsPerPeriodRow
	for rows.Next() {
		var i GetMetricStatsPerPeriodRow
		if err := rows.Scan(
			&i.Source,
			&i.Name,
			&i.Ts,
			&i.Avg,
			&i.Max,
			&i.Min,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const insertMetric = `-- name: InsertMetric :one
INSERT INTO public.metrics (account_id, ts, source, name, target, value, inserted_at)
VALUES ($1, $2, $3, $4, $5, $6, NOW())
RETURNING source, ts, inserted_at, name, target, value, ip_address, account_id
`

type InsertMetricParams struct {
	AccountID int64           `json:"account_id"`
	Ts        sql.NullTime    `json:"ts"`
	Source    string          `json:"source"`
	Name      string          `json:"name"`
	Target    string          `json:"target"`
	Value     sql.NullFloat64 `json:"value"`
}

func (q *Queries) InsertMetric(ctx context.Context, arg InsertMetricParams) (Metric, error) {
	row := q.db.QueryRowContext(ctx, insertMetric,
		arg.AccountID,
		arg.Ts,
		arg.Source,
		arg.Name,
		arg.Target,
		arg.Value,
	)
	var i Metric
	err := row.Scan(
		&i.Source,
		&i.Ts,
		&i.InsertedAt,
		&i.Name,
		&i.Target,
		&i.Value,
		&i.IpAddress,
		&i.AccountID,
	)
	return i, err
}
