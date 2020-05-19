package tracer

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"github.com/lib/pq"
	"github.com/ngrok/sqlmw"
	"go.opentelemetry.io/otel/api/kv"
	"go.opentelemetry.io/otel/api/trace"
	"runtime"
	"strings"
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
	ctx, dbSpan := span.Tracer().Start(ctx, "postgres-helpful-query-name", trace.WithSpanKind(trace.SpanKindClient))
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
	queryLabel := getCallerName()
	if queryLabel == "" {
		queryLabel = "postgres-query"
	}
	ctx, dbSpan := span.Tracer().Start(ctx, queryLabel, trace.WithSpanKind(trace.SpanKindClient))
	dbSpan.SetAttributes(
		kv.String("db.kind", "sql"),
		kv.String("db.instance", ""),
		kv.String("db.statement", query),
		kv.String("db.user", "app"),
	)
	defer dbSpan.End()
	return conn.QueryContext(ctx, query, args)
}

var dbLibs = []string{
	"com/ngrok/sqlmw",
	"database/sql",
	"sqlInterceptor",
}

// we want to introspect the first caller that is outside of our tracing/database/orm libs
func getCallerName() string {
	pc := make([]uintptr, 10)
	frameCount := runtime.Callers(3, pc)
	frames := runtime.CallersFrames(pc[:frameCount])
	frame, more := frames.Next()
	// find the first calling function outside of database/sql and ngrok
	for more {
		if frame.Func == nil {
			return ""
		}
		dbLib := false
		for i := 0; i < len(dbLibs) && !dbLib; i++ {
			if strings.Contains(frame.Func.Name(), dbLibs[i]) {
				dbLib = true
			}
		}
		if !dbLib {
			return frame.Func.Name()
		}
		frame, more = frames.Next()
	}
	return ""
}
