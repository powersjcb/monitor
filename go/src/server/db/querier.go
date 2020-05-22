// Code generated by sqlc. DO NOT EDIT.

package db

import (
	"context"
)

type Querier interface {
	GetAccountIDForAPIKey(ctx context.Context, apiKey string) (int64, error)
	GetMetricForSource(ctx context.Context, source string) ([]Metric, error)
	GetMetricStatsPerPeriod(ctx context.Context, seconds int32) ([]GetMetricStatsPerPeriodRow, error)
	GetMetrics(ctx context.Context) ([]string, error)
	InsertAccount(ctx context.Context, arg InsertAccountParams) (Account, error)
	InsertMetric(ctx context.Context, arg InsertMetricParams) (Metric, error)
}

var _ Querier = (*Queries)(nil)
