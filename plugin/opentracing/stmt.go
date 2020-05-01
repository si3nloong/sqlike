package opentracing

import (
	"context"
	"database/sql/driver"

	"github.com/opentracing/opentracing-go/ext"
)

// StmtExecContext :
func (ot *OpenTracingInterceptor) StmtExecContext(ctx context.Context, conn driver.StmtExecContext, query string, args []driver.NamedValue) (driver.Result, error) {
	span := ot.StartSpan(ctx, "exec")
	defer span.Finish()
	ext.DBStatement.Set(span, query)
	rows, err := conn.ExecContext(ctx, args)
	if err != nil {
		ext.LogError(span, err)
		return nil, err
	}
	return rows, nil
}

// StmtQueryContext :
func (ot *OpenTracingInterceptor) StmtQueryContext(ctx context.Context, conn driver.StmtQueryContext, query string, args []driver.NamedValue) (driver.Rows, error) {
	span := ot.StartSpan(ctx, "query_context")
	defer span.Finish()
	ext.DBStatement.Set(span, query)
	rows, err := conn.QueryContext(ctx, args)
	if err != nil {
		ext.LogError(span, err)
		return nil, err
	}
	return rows, nil
}

// StmtClose :
func (ot *OpenTracingInterceptor) StmtClose(ctx context.Context, conn driver.Stmt) error {
	span := ot.StartSpan(ctx, "close")
	defer span.Finish()
	if err := conn.Close(); err != nil {
		ext.LogError(span, err)
		return err
	}
	return nil
}
