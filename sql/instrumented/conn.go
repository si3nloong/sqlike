package instrumented

import (
	"context"
	"database/sql/driver"
)

type Conn interface {
	driver.Conn
	driver.Pinger
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.QueryerContext
}

type wrappedConn struct {
	itpr Interceptor
	conn Conn
}

var _ Conn = (*wrappedConn)(nil)

// Begin :
func (w wrappedConn) Begin() (driver.Tx, error) {
	tx, err := w.conn.Begin()
	if err != nil {
		return nil, err
	}
	return wrappedTx{ctx: context.Background(), itpr: w.itpr, tx: tx}, nil
}

// BeginTx :
func (w wrappedConn) BeginTx(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
	tx, err := w.itpr.ConnBeginTx(ctx, w.conn, opts)
	if err != nil {
		return nil, err
	}
	return wrappedTx{ctx: ctx, itpr: w.itpr, tx: tx}, nil
}

// Ping :
func (w wrappedConn) Ping(ctx context.Context) error {
	return w.itpr.ConnPing(ctx, w.conn)
}

// Prepare :
func (w wrappedConn) Prepare(query string) (driver.Stmt, error) {
	ctx := context.Background()
	stmt, err := w.itpr.ConnPrepareContext(ctx, w.conn, query)
	if err != nil {
		return nil, err
	}
	x, ok := stmt.(Stmt)
	if !ok {
		return nil, driver.ErrBadConn
	}
	return wrappedStmt{ctx: ctx, itpr: w.itpr, query: query, stmt: x}, nil
}

// PrepareContext :
func (w wrappedConn) PrepareContext(ctx context.Context, query string) (driver.Stmt, error) {
	stmt, err := w.itpr.ConnPrepareContext(ctx, w.conn, query)
	if err != nil {
		return nil, err
	}
	x, ok := stmt.(Stmt)
	if !ok {
		return nil, driver.ErrBadConn
	}
	return wrappedStmt{ctx: ctx, itpr: w.itpr, query: query, stmt: x}, nil
}

// ExecContext :
func (w wrappedConn) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (res driver.Result, err error) {
	result, err := w.itpr.ConnExecContext(ctx, w.conn, query, args)
	if err != nil {
		return nil, err
	}
	return wrappedResult{ctx: ctx, itpr: w.itpr, result: result}, nil
}

// QueryContext :
func (w wrappedConn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (res driver.Rows, err error) {
	rows, err := w.itpr.ConnQueryContext(ctx, w.conn, query, args)
	if err != nil {
		return nil, err
	}
	x, ok := rows.(Rows)
	if !ok {
		return nil, driver.ErrSkip
	}
	return wrappedRows{ctx: ctx, itpr: w.itpr, rows: x}, nil
}

// Close :
func (w wrappedConn) Close() error {
	return w.conn.Close()
}
