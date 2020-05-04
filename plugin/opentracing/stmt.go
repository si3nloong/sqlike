package opentracing

import (
	"context"
	"database/sql/driver"

	"github.com/opentracing/opentracing-go"
)

// StmtExecContext :
func (ot *OpenTracingInterceptor) StmtExecContext(ctx context.Context, conn driver.StmtExecContext, query string, args []driver.NamedValue) (result driver.Result, err error) {
	if ot.opts.Exec {
		var span opentracing.Span
		span, ctx = ot.MaybeStartSpanFromContext(ctx, "exec")
		ot.logQuery(span, query)
		ot.logArgs(span, args)
		defer func() {
			ot.logError(span, err)
			span.Finish()
		}()
	}
	result, err = conn.ExecContext(ctx, args)
	return
}

// StmtQueryContext :
func (ot *OpenTracingInterceptor) StmtQueryContext(ctx context.Context, conn driver.StmtQueryContext, query string, args []driver.NamedValue) (rows driver.Rows, err error) {
	if ot.opts.Query {
		var span opentracing.Span
		span, ctx = ot.MaybeStartSpanFromContext(ctx, "query")
		ot.logQuery(span, query)
		ot.logArgs(span, args)
		defer func() {
			ot.logError(span, err)
			span.Finish()
		}()
	}
	rows, err = conn.QueryContext(ctx, args)
	return
}

// StmtClose :
func (ot *OpenTracingInterceptor) StmtClose(ctx context.Context, conn driver.Stmt) (err error) {
	if ot.opts.RowsClose {
		var span opentracing.Span
		span, ctx = ot.MaybeStartSpanFromContext(ctx, "close")
		defer func() {
			ot.logError(span, err)
			span.Finish()
		}()
	}
	err = conn.Close()
	return nil
}
