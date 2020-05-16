package tracer

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"github.com/lib/pq"
	"github.com/ngrok/sqlmw"
	"go.opentelemetry.io/otel/api/kv"
	"go.opentelemetry.io/otel/api/trace"
)

func OpenDB(dsn string) (*sql.DB, error) {
	sql.Register("postgres-tracing", sqlmw.Driver(&pq.Driver{}, new(sqlInterceptor)))
	db, err := sql.Open("postgres-tracing", dsn)
	return db, err
}

type sqlInterceptor struct {
	sqlmw.NullInterceptor
}

// ref: https://github.com/open-telemetry/opentelemetry-specification/blob/master/specification/trace/semantic_conventions/database.md
func (sqlInterceptor) ConnExecContext(ctx context.Context, conn driver.ExecerContext, query string, args []driver.NamedValue) (driver.Result, error) {
	span := trace.SpanFromContext(ctx)
	// assumes query does not contain any values
	ctx, dbSpan := span.Tracer().Start(ctx, "postgres-helpfull-query-name", trace.WithSpanKind(trace.SpanKindClient))
	dbSpan.SetAttributes(
		kv.String("db.kind", "sql"),
		kv.String("db.instance", ""),
		kv.String("db.statement", query),
		kv.String("db.user", ""),
	)
	defer dbSpan.End()
	return conn.ExecContext(ctx, query, args)
}

func (sqlInterceptor) ConnQueryContext(ctx context.Context, conn driver.QueryerContext, query string, args []driver.NamedValue) (driver.Rows, error) {
	span := trace.SpanFromContext(ctx)
	// assumes query does not contain any values
	ctx, dbSpan := span.Tracer().Start(ctx, "postgres-helpfull-query-name", trace.WithSpanKind(trace.SpanKindClient))
	dbSpan.SetAttributes(
		kv.String("db.kind", "sql"),
		kv.String("db.instance", ""),
		kv.String("db.statement", query),
		kv.String("db.user", "app"),
	)
	defer dbSpan.End()
	return conn.QueryContext(ctx, query, args)
}