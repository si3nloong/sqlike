package instrumented

import (
	"context"
	"database/sql/driver"
)

type WrappedConn struct {
	conn driver.Conn
}

var (
	_ driver.Conn        = (*WrappedConn)(nil)
	_ driver.ConnBeginTx = (*WrappedConn)(nil)
)

func (w WrappedConn) Begin() (driver.Tx, error) {
	tx, err := w.conn.Begin()
	if err != nil {
		return nil, err
	}
	return WrappedTx{tx: tx}, nil
}

func (w WrappedConn) BeginTx(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
	x, ok := w.conn.(driver.ConnBeginTx)
	if ok {
		tx, err := x.BeginTx(ctx, opts)
		if err != nil {
			return nil, err
		}
		return WrappedTx{tx: tx}, nil
	}
	return w.Begin()
}

func (w WrappedConn) Prepare(query string) (driver.Stmt, error) {
	stmt, err := w.conn.Prepare(query)
	if err != nil {
		return nil, err
	}
	return WrappedStmt{stmt: stmt}, nil
}

func (w WrappedConn) Close() error {
	return w.conn.Close()
}
