package db

import (
	"context"

	"github.com/si3nloong/sqlike/sql"
)

// ColumnDataTyper :
type ColumnDataTyper interface {
	ColumnDataType(ctx context.Context) *sql.Column
}
