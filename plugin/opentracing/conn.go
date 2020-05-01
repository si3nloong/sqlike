package opentracing

import (
	"context"
	"database/sql/driver"

	"github.com/opentracing/opentracing-go/ext"
)

// ConnPing :
func (ot *OpenTracingInterceptor) ConnPing(ctx context.Context, conn driver.Pinger) error {
	span := ot.StartSpan(ctx, "conn_ping")
	defer span.Finish()
	if err := conn.Ping(ctx); err != nil {
		ext.LogError(span, err)
		return err
	}
	return nil
}

// ConnPing :
func (ot *OpenTracingInterceptor) ConnBeginTx(ctx context.Context, conn driver.ConnBeginTx, opts driver.TxOptions) (driver.Tx, error) {
	span := ot.StartSpan(ctx, "conn_begin_transaction")
	defer span.Finish()
	tx, err := conn.BeginTx(ctx, opts)
	if err != nil {
		ext.LogError(span, err)
		return nil, err
	}
	return tx, nil
}

// ConnPrepareContext :
func (ot *OpenTracingInterceptor) ConnPrepareContext(ctx context.Context, conn driver.ConnPrepareContext, query string) (driver.Stmt, error) {
	span := ot.StartSpan(ctx, "conn_prepare")
	defer span.Finish()
	ext.DBStatement.Set(span, query)
	stmt, err := conn.PrepareContext(ctx, query)
	if err != nil {
		ext.LogError(span, err)
		return nil, err
	}
	return stmt, nil
}

// ConnExecContext :
func (ot *OpenTracingInterceptor) ConnExecContext(ctx context.Context, conn driver.ExecerContext, query string, args []driver.NamedValue) (driver.Result, error) {
	span := ot.StartSpan(ctx, "conn_exec")
	defer span.Finish()
	ext.DBStatement.Set(span, query)
	result, err := conn.ExecContext(ctx, query, args)
	if err != nil {
		ext.LogError(span, err)
		return nil, err
	}
	return result, nil
}

func (ot *OpenTracingInterceptor) ConnQueryContext(ctx context.Context, conn driver.QueryerContext, query string, args []driver.NamedValue) (driver.Rows, error) {
	span := ot.StartSpan(ctx, "conn_query")
	defer span.Finish()
	ext.DBStatement.Set(span, query)
	rows, err := conn.QueryContext(ctx, query, args)
	if err != nil {
		ext.LogError(span, err)
		return nil, err
	}
	return rows, nil
}
