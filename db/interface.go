package db

import (
	"context"

	"github.com/si3nloong/sqlike/sql"
)

// SQLValuer :
type SQLValuer interface {
	SQLValue(ctx context.Context) (interface{}, error)
}

// SQLScanner :
type SQLScanner interface {
	SQLScan(ctx context.Context, val interface{}) error
}

// ColumnDataTypeImplementer :
type ColumnDataTypeImplementer interface {
	ColumnDataType(ctx context.Context) *sql.Column
}
