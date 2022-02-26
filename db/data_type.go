package db

import (
	"context"

	"github.com/si3nloong/sqlike/v2/sql"
)

// ColumnDataTyper :
type ColumnDataTyper interface {
	ColumnDataType(ctx context.Context) *sql.Column
}
