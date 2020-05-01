package instrumented

import (
	"context"
	"database/sql/driver"
)

type Stmt interface {
	driver.Stmt
	driver.StmtExecContext
	driver.StmtQueryContext
}

type wrappedStmt struct {
	ctx   context.Context
	itpr  Interceptor
	query string
	stmt  Stmt
}

var _ Stmt = (*wrappedStmt)(nil)

// Exec :
func (w wrappedStmt) Exec(args []driver.Value) (driver.Result, error) {
	result, err := w.stmt.Exec(args)
	if err != nil {
		return nil, err
	}
	return wrappedResult{ctx: w.ctx, itpr: w.itpr, result: result}, nil
}

// ExecContext :
func (w wrappedStmt) ExecContext(ctx context.Context, args []driver.NamedValue) (driver.Result, error) {
	result, err := w.itpr.StmtExecContext(ctx, w.stmt, w.query, args)
	if err != nil {
		return nil, err
	}
	return wrappedResult{ctx: ctx, itpr: w.itpr, result: result}, nil
}

// Query :
func (w wrappedStmt) Query(args []driver.Value) (driver.Rows, error) {
	rows, err := w.stmt.Query(args)
	if err != nil {
		return nil, err
	}
	x, ok := rows.(Rows)
	if !ok {
		return nil, driver.ErrSkip
	}
	return wrappedRows{ctx: w.ctx, itpr: w.itpr, rows: x}, nil
}

// QueryContext :
func (w wrappedStmt) QueryContext(ctx context.Context, args []driver.NamedValue) (driver.Rows, error) {
	rows, err := w.itpr.StmtQueryContext(ctx, w.stmt, w.query, args)
	if err != nil {
		return nil, err
	}
	x, ok := rows.(Rows)
	if !ok {
		return nil, driver.ErrSkip
	}
	return wrappedRows{ctx: ctx, itpr: w.itpr, rows: x}, nil
}

// NumInput :
func (w wrappedStmt) NumInput() int {
	return w.stmt.NumInput()
}

// Close :
func (w wrappedStmt) Close() error {
	return w.itpr.StmtClose(w.ctx, w.stmt)
}
