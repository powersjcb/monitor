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

type sqlInterceptor struct {}

// ref: https://github.com/open-telemetry/opentelemetry-specification/blob/master/specification/trace/semantic_conventions/database.md
const dbKind = "db.kind"
const dbInstance = "db.instance"
const dbStatement = "db.statement"
const dbUser = "db.user"

const defaultQueryLabel = "postgres"
const defaultDBUser = "app"
const dbKindValue = "sql"

func (sqlInterceptor) ConnBeginTx(ctx context.Context, conn driver.ConnBeginTx, txOpts driver.TxOptions) (driver.Tx, error) {
	span := trace.SpanFromContext(ctx)
	queryLabel := getCallerName()
	if queryLabel == "" {
		queryLabel = defaultQueryLabel
	}
	ctx, dbSpan := span.Tracer().Start(ctx, queryLabel, trace.WithSpanKind(trace.SpanKindClient))
	dbSpan.SetAttributes(
		kv.String(dbKind, dbKindValue),
		kv.String(dbInstance, ""),
		kv.String(dbStatement, "BEGIN"),
		kv.String(dbUser, defaultDBUser),
	)
	defer dbSpan.End()
	return conn.BeginTx(ctx, txOpts)
}

func (sqlInterceptor) ConnPrepareContext(ctx context.Context, conn driver.ConnPrepareContext, query string) (driver.Stmt, error) {
	span := trace.SpanFromContext(ctx)
	queryLabel := getCallerName()
	if queryLabel == "" {
		queryLabel = defaultQueryLabel
	}
	ctx, dbSpan := span.Tracer().Start(ctx, queryLabel, trace.WithSpanKind(trace.SpanKindClient))
	dbSpan.SetAttributes(
		kv.String(dbKind, dbKindValue),
		kv.String(dbInstance, ""),
		kv.String(dbStatement, query), // todo: determine what data is provided in query for preparing a statement
		kv.String(dbUser, defaultDBUser),
	)
	defer dbSpan.End()
	return conn.PrepareContext(ctx, query)
}

func (sqlInterceptor) ConnPing(ctx context.Context, conn driver.Pinger) error {
	span := trace.SpanFromContext(ctx)
	queryLabel := getCallerName()
	if queryLabel == "" {
		queryLabel = defaultQueryLabel
	}
	ctx, dbSpan := span.Tracer().Start(ctx, queryLabel, trace.WithSpanKind(trace.SpanKindClient))
	dbSpan.SetAttributes(
		kv.String(dbKind, dbKindValue),
		kv.String(dbInstance, ""),
		kv.String(dbStatement, "PING"),
		kv.String(dbUser, defaultDBUser),
	)
	defer dbSpan.End()
	return conn.Ping(ctx)
}

func (sqlInterceptor) ConnExecContext(ctx context.Context, conn driver.ExecerContext, query string, args []driver.NamedValue) (driver.Result, error) {
	span := trace.SpanFromContext(ctx)
	queryLabel := getCallerName()
	if queryLabel == "" {
		queryLabel = defaultQueryLabel
	}
	ctx, dbSpan := span.Tracer().Start(ctx, queryLabel, trace.WithSpanKind(trace.SpanKindClient))
	dbSpan.SetAttributes(
		kv.String(dbKind, dbKindValue),
		kv.String(dbInstance, ""),
		kv.String(dbStatement, query),
		kv.String(dbUser, defaultDBUser),
	)
	defer dbSpan.End()
	return conn.ExecContext(ctx, query, args)
}

func (sqlInterceptor) ConnQueryContext(ctx context.Context, conn driver.QueryerContext, query string, args []driver.NamedValue) (driver.Rows, error) {
	span := trace.SpanFromContext(ctx)
	queryLabel := getCallerName()
	if queryLabel == "" {
		queryLabel = defaultQueryLabel
	}
	ctx, dbSpan := span.Tracer().Start(ctx, queryLabel, trace.WithSpanKind(trace.SpanKindClient))
	dbSpan.SetAttributes(
		kv.String(dbKind, dbKindValue),
		kv.String(dbInstance, ""),
		kv.String(dbStatement, query),
		kv.String(dbUser, defaultDBUser),
	)
	defer dbSpan.End()
	return conn.QueryContext(ctx, query, args)
}

// gets a datatbase conn for a single goroutine
func (sqlInterceptor) ConnectorConnect(ctx context.Context, connect driver.Connector) (driver.Conn, error) {
	span := trace.SpanFromContext(ctx)
	queryLabel := getCallerName()
	if queryLabel == "" {
		queryLabel = defaultQueryLabel
	}
	ctx, dbSpan := span.Tracer().Start(ctx, queryLabel, trace.WithSpanKind(trace.SpanKindClient))
	dbSpan.SetAttributes(
		kv.String(dbKind, dbKindValue),
		kv.String(dbInstance, ""),
		kv.String(dbStatement, "CONNECT"),
		kv.String(dbUser, defaultDBUser),
	)
	defer dbSpan.End()
	return connect.Connect(ctx)
}

func (sqlInterceptor) ResultLastInsertId(res driver.Result) (int64, error) {
	return res.LastInsertId()
}

func (sqlInterceptor) ResultRowsAffected(res driver.Result) (int64, error) {
	return res.RowsAffected()
}

func (sqlInterceptor) RowsNext(_ context.Context, rows driver.Rows, dest []driver.Value) error {
	return rows.Next(dest)
}

func (sqlInterceptor) StmtExecContext(ctx context.Context, stmt driver.StmtExecContext, query string, args []driver.NamedValue) (driver.Result, error) {
	span := trace.SpanFromContext(ctx)
	queryLabel := getCallerName()
	if queryLabel == "" {
		queryLabel = defaultQueryLabel
	}
	ctx, dbSpan := span.Tracer().Start(ctx, queryLabel, trace.WithSpanKind(trace.SpanKindClient))
	dbSpan.SetAttributes(
		kv.String(dbKind, dbKindValue),
		kv.String(dbInstance, ""),
		kv.String(dbStatement, query),
		kv.String(dbUser, defaultDBUser),
	)
	defer dbSpan.End()
	return stmt.ExecContext(ctx, args)
}

func (sqlInterceptor) StmtQueryContext(ctx context.Context, stmt driver.StmtQueryContext, query string, args []driver.NamedValue) (driver.Rows, error) {
	span := trace.SpanFromContext(ctx)
	queryLabel := getCallerName()
	if queryLabel == "" {
		queryLabel = defaultQueryLabel
	}
	ctx, dbSpan := span.Tracer().Start(ctx, queryLabel, trace.WithSpanKind(trace.SpanKindClient))
	dbSpan.SetAttributes(
		kv.String(dbKind, dbKindValue),
		kv.String(dbInstance, ""),
		kv.String(dbStatement, query),
		kv.String(dbUser, defaultDBUser),
	)
	defer dbSpan.End()
	return stmt.QueryContext(ctx, args)
}

func (sqlInterceptor) StmtClose(ctx context.Context, stmt driver.Stmt) error {
	span := trace.SpanFromContext(ctx)
	queryLabel := getCallerName()
	if queryLabel == "" {
		queryLabel = defaultQueryLabel
	}
	_, dbSpan := span.Tracer().Start(ctx, queryLabel, trace.WithSpanKind(trace.SpanKindClient))
	dbSpan.SetAttributes(
		kv.String(dbKind, dbKindValue),
		kv.String(dbInstance, ""),
		kv.String(dbStatement, "CLOSE PREPARED STATEMENT"),
		kv.String(dbUser, defaultDBUser),
	)
	defer dbSpan.End()
	return stmt.Close()
}

func (sqlInterceptor) TxCommit(ctx context.Context, tx driver.Tx) error {
	span := trace.SpanFromContext(ctx)
	_, dbSpan := span.Tracer().Start(ctx, defaultQueryLabel, trace.WithSpanKind(trace.SpanKindClient))
	dbSpan.SetAttributes(
		kv.String(dbKind, dbKindValue),
		kv.String(dbInstance, ""),
		kv.String(dbStatement, "COMMIT"),
		kv.String(dbUser, defaultDBUser),
	)
	defer dbSpan.End()
	return tx.Commit()
}

func (sqlInterceptor) TxRollback(ctx context.Context, tx driver.Tx) error {
	span := trace.SpanFromContext(ctx)
	_, dbSpan := span.Tracer().Start(ctx, defaultQueryLabel, trace.WithSpanKind(trace.SpanKindClient))
	dbSpan.SetAttributes(
		kv.String(dbKind, dbKindValue),
		kv.String(dbInstance, ""),
		kv.String(dbStatement, "ROLLBACK"),
		kv.String(dbUser, defaultDBUser),
	)
	defer dbSpan.End()
	return tx.Rollback()
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
