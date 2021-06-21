package db

import (
	"context"

	"github.com/si3nloong/sqlike/sqlike/columns"
)

// SQLValuer :
type SQLValuer interface {
	SQLValue(ctx context.Context) (interface{}, error)
}

// ColumnDataTypeImplementer :
type ColumnDataTypeImplementer interface {
	ColumnDataType(ctx context.Context) *columns.Column
}
