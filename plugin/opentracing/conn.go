package opentracing

import (
	"context"
	"database/sql/driver"

	"github.com/opentracing/opentracing-go/ext"
)

// ConnPing :
func (ot *OpenTracingInterceptor) ConnPing(ctx context.Context, conn driver.Pinger) (err error) {
	if ot.opts.Ping {
		span := ot.StartSpan(ctx, "ping")
		defer func() {
			if err != nil {
				ext.LogError(span, err)
			}
			span.Finish()
		}()
	}
	err = conn.Ping(ctx)
	return
}

// ConnPing :
func (ot *OpenTracingInterceptor) ConnBeginTx(ctx context.Context, conn driver.ConnBeginTx, opts driver.TxOptions) (tx driver.Tx, err error) {
	if ot.opts.BeginTx {
		span := ot.StartSpan(ctx, "begin_tx")
		defer func() {
			if err != nil {
				ext.LogError(span, err)
			}
			span.Finish()
		}()
	}
	tx, err = conn.BeginTx(ctx, opts)
	return
}

// ConnPrepareContext :
func (ot *OpenTracingInterceptor) ConnPrepareContext(ctx context.Context, conn driver.ConnPrepareContext, query string) (stmt driver.Stmt, err error) {
	if ot.opts.Prepare {
		span := ot.StartSpan(ctx, "prepare")
		ext.DBStatement.Set(span, query)
		defer func() {
			if err != nil {
				ext.LogError(span, err)
			}
			span.Finish()
		}()
	}
	stmt, err = conn.PrepareContext(ctx, query)
	return
}

// ConnExecContext :
func (ot *OpenTracingInterceptor) ConnExecContext(ctx context.Context, conn driver.ExecerContext, query string, args []driver.NamedValue) (result driver.Result, err error) {
	if ot.opts.Exec {
		span := ot.StartSpan(ctx, "exec")
		ext.DBStatement.Set(span, query)
		if ot.opts.Args {
			logArgs(span, args)
		}
		defer func() {
			if err != nil {
				ext.LogError(span, err)
			}
			span.Finish()
		}()
	}
	result, err = conn.ExecContext(ctx, query, args)
	return
}

func (ot *OpenTracingInterceptor) ConnQueryContext(ctx context.Context, conn driver.QueryerContext, query string, args []driver.NamedValue) (rows driver.Rows, err error) {
	if ot.opts.Query {
		span := ot.StartSpan(ctx, "query")
		ext.DBStatement.Set(span, query)
		if ot.opts.Args {
			logArgs(span, args)
		}
		defer func() {
			if err != nil {
				ext.LogError(span, err)
			}
			span.Finish()
		}()
	}
	rows, err = conn.QueryContext(ctx, query, args)
	return
}
