package opentracing

import (
	"context"
	"database/sql/driver"

	"github.com/opentracing/opentracing-go"
)

// ConnPing :
func (ot *OpenTracingInterceptor) ConnPing(ctx context.Context, conn driver.Pinger) (err error) {
	if ot.opts.Ping {
		var span opentracing.Span
		span, ctx = ot.MaybeStartSpanFromContext(ctx, "ping")
		defer func() {
			ot.logError(span, err)
			span.Finish()
		}()
	}
	err = conn.Ping(ctx)
	return
}

// ConnBeginTx :
func (ot *OpenTracingInterceptor) ConnBeginTx(ctx context.Context, conn driver.ConnBeginTx, opts driver.TxOptions) (tx driver.Tx, err error) {
	if ot.opts.BeginTx {
		var span opentracing.Span
		span, ctx = ot.MaybeStartSpanFromContext(ctx, "begin_tx")
		defer func() {
			ot.logError(span, err)
			span.Finish()
		}()
	}
	tx, err = conn.BeginTx(ctx, opts)
	return
}

// ConnPrepareContext :
func (ot *OpenTracingInterceptor) ConnPrepareContext(ctx context.Context, conn driver.ConnPrepareContext, query string) (stmt driver.Stmt, err error) {
	if ot.opts.Prepare {
		var span opentracing.Span
		span, ctx = ot.MaybeStartSpanFromContext(ctx, "prepare")
		// ext.DBStatement.Set(span, query)
		ot.logQuery(span, query)
		defer func() {
			ot.logError(span, err)
			span.Finish()
		}()
	}
	stmt, err = conn.PrepareContext(ctx, query)
	return
}
