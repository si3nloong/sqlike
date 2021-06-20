package db

import (
	"context"

	"github.com/si3nloong/sqlike/sqlike/columns"
)

// SqlValueConverter :
type SqlValueConverter interface {
	SqlValue(ctx context.Context) (interface{}, error)
}

// ColumnDataTypeImplementer :
type ColumnDataTypeImplementer interface {
	ColumnDataType(ctx context.Context) *columns.Column
}
