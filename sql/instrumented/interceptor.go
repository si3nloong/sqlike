package instrumented

import (
	"context"
	"database/sql/driver"
)

// Interceptor :
type Interceptor interface {
	// Connection interceptors
	ConnPing(context.Context, driver.Pinger) error
	ConnBeginTx(context.Context, driver.ConnBeginTx, driver.TxOptions) (driver.Tx, error)
	ConnPrepareContext(context.Context, driver.ConnPrepareContext, string) (driver.Stmt, error)
	ConnExecContext(context.Context, driver.ExecerContext, string, []driver.NamedValue) (driver.Result, error)
	ConnQueryContext(context.Context, driver.QueryerContext, string, []driver.NamedValue) (driver.Rows, error)

	// Connector interceptors
	// ConnectorConnect(context.Context, driver.Connector) (driver.Conn, error)

	// Results interceptors
	ResultLastInsertId(context.Context, driver.Result) (int64, error)
	ResultRowsAffected(context.Context, driver.Result) (int64, error)

	// Rows interceptors
	RowsNext(context.Context, driver.Rows, []driver.Value) error

	// Stmt interceptors
	StmtExecContext(context.Context, driver.StmtExecContext, string, []driver.NamedValue) (driver.Result, error)
	StmtQueryContext(context.Context, driver.StmtQueryContext, string, []driver.NamedValue) (driver.Rows, error)
	StmtClose(context.Context, driver.Stmt) error

	// Tx interceptors
	TxCommit(context.Context, driver.Tx) error
	TxRollback(context.Context, driver.Tx) error
}

var _ Interceptor = NullInterceptor{}

// NullInterceptor is a complete passthrough interceptor that implements every method of the Interceptor
// interface and performs no additional logic. Users should Embed it in their own interceptor so that they
// only need to define the specific functions they are interested in intercepting.
type NullInterceptor struct{}

// ConnBeginTx :
func (NullInterceptor) ConnBeginTx(ctx context.Context, conn driver.ConnBeginTx, txOpts driver.TxOptions) (driver.Tx, error) {
	return conn.BeginTx(ctx, txOpts)
}

// ConnPrepareContext :
func (NullInterceptor) ConnPrepareContext(ctx context.Context, conn driver.ConnPrepareContext, query string) (driver.Stmt, error) {
	return conn.PrepareContext(ctx, query)
}

// ConnPing :
func (NullInterceptor) ConnPing(ctx context.Context, conn driver.Pinger) error {
	return conn.Ping(ctx)
}

// ConnExecContext :
func (NullInterceptor) ConnExecContext(ctx context.Context, conn driver.ExecerContext, query string, args []driver.NamedValue) (driver.Result, error) {
	return conn.ExecContext(ctx, query, args)
}

// ConnQueryContext :
func (NullInterceptor) ConnQueryContext(ctx context.Context, conn driver.QueryerContext, query string, args []driver.NamedValue) (driver.Rows, error) {
	return conn.QueryContext(ctx, query, args)
}

// func (NullInterceptor) ConnectorConnect(ctx context.Context, connect driver.Connector) (driver.Conn, error) {
// 	return connect.Connect(ctx)
// }

// ResultLastInsertId :
func (NullInterceptor) ResultLastInsertId(ctx context.Context, res driver.Result) (int64, error) {
	return res.LastInsertId()
}

// ResultRowsAffected :
func (NullInterceptor) ResultRowsAffected(ctx context.Context, res driver.Result) (int64, error) {
	return res.RowsAffected()
}

// RowsNext :
func (NullInterceptor) RowsNext(ctx context.Context, rows driver.Rows, dest []driver.Value) error {
	return rows.Next(dest)
}

// StmtExecContext :
func (NullInterceptor) StmtExecContext(ctx context.Context, stmt driver.StmtExecContext, _ string, args []driver.NamedValue) (driver.Result, error) {
	return stmt.ExecContext(ctx, args)
}

// StmtQueryContext :
func (NullInterceptor) StmtQueryContext(ctx context.Context, stmt driver.StmtQueryContext, _ string, args []driver.NamedValue) (driver.Rows, error) {
	return stmt.QueryContext(ctx, args)
}

// StmtClose :
func (NullInterceptor) StmtClose(ctx context.Context, stmt driver.Stmt) error {
	return stmt.Close()
}

// TxCommit :
func (NullInterceptor) TxCommit(ctx context.Context, tx driver.Tx) error {
	return tx.Commit()
}

// TxRollback :
func (NullInterceptor) TxRollback(ctx context.Context, tx driver.Tx) error {
	return tx.Rollback()
}
