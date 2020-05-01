package instrumented

import (
	"context"
	"database/sql/driver"
	"reflect"
)

// Rows :
type Rows interface {
	driver.RowsNextResultSet
	ColumnTypeScanType(index int) reflect.Type
}

type wrappedRows struct {
	ctx  context.Context
	itpr Interceptor
	rows Rows
}

var _ Rows = (*wrappedRows)(nil)

// Columns :
func (w wrappedRows) Columns() []string {
	return w.rows.Columns()
}

// Next :
func (w wrappedRows) Next(dest []driver.Value) error {
	return w.itpr.RowsNext(w.ctx, w.rows, dest)
}

// HasNextResultSet :
func (w wrappedRows) HasNextResultSet() bool {
	return w.rows.HasNextResultSet()
}

// NextResultSet :
func (w wrappedRows) NextResultSet() error {
	return w.rows.NextResultSet()
}

// ColumnTypeScanType :
func (w wrappedRows) ColumnTypeScanType(index int) reflect.Type {
	return w.rows.ColumnTypeScanType(index)
}

// Close :
func (w wrappedRows) Close() error {
	return w.rows.Close()
}
