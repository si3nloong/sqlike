package instrumented

import (
	"context"
	"database/sql/driver"
)

type WrappedStmt struct {
	stmt driver.Stmt
}

var (
	_ driver.Stmt             = (*WrappedStmt)(nil)
	_ driver.StmtExecContext  = (*WrappedStmt)(nil)
	_ driver.StmtQueryContext = (*WrappedStmt)(nil)
)

// Exec :
func (w WrappedStmt) Exec(args []driver.Value) (driver.Result, error) {
	result, err := w.stmt.Exec(args)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// ExecContext :
func (w WrappedStmt) ExecContext(ctx context.Context, args []driver.NamedValue) (driver.Result, error) {
	x, ok := w.stmt.(driver.StmtExecContext)
	if ok {
		result, err := x.ExecContext(ctx, args)
		if err != nil {
			return nil, err
		}
		return result, nil
	}
	// TODO: convert NamedValue to Value
	return w.Exec(nil)
}

// QueryContext :
func (w WrappedStmt) QueryContext(ctx context.Context, args []driver.NamedValue) (driver.Rows, error) {
	x, ok := w.stmt.(driver.StmtQueryContext)
	if ok {
		result, err := x.QueryContext(ctx, args)
		if err != nil {
			return nil, err
		}
		return result, nil
	}
	// TODO: convert NamedValue to Value
	return w.Query(nil)
}

// Query :
func (w WrappedStmt) Query(args []driver.Value) (driver.Rows, error) {
	return w.stmt.Query(args)
}

// NumInput :
func (w WrappedStmt) NumInput() int {
	return w.stmt.NumInput()
}

// Close :
func (w WrappedStmt) Close() error {
	return w.stmt.Close()
}
