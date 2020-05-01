package instrumented

import (
	"context"
	"database/sql/driver"
)

type wrappedRows struct {
	ctx  context.Context
	itpr Interceptor
	rows driver.Rows
}

// Columns :
func (w wrappedRows) Columns() []string {
	return w.rows.Columns()
}

// Next :
func (w wrappedRows) Next(dest []driver.Value) error {
	return w.itpr.RowsNext(w.ctx, w.rows, dest)
}

// Close :
func (w wrappedRows) Close() error {
	return w.rows.Close()
}
