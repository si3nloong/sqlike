package sql

import (
	"fmt"

	"github.com/si3nloong/sqlike/v2/internal/primitive"
)

// UpdateStmt :
type UpdateStmt struct {
	Database   string
	Table      string
	Conditions primitive.Group
	Values     []primitive.KV
	Sorts      []any
	RowCount   uint
}

// Update :
func Update[T ~string | primitive.Pair](src T) *UpdateStmt {
	stmt := new(UpdateStmt)
	switch vi := any(src).(type) {
	case primitive.Pair:
		stmt.Database = vi[0]
		stmt.Table = vi[1]
	case string:
		stmt.Table = vi
	default:
		stmt.Table = fmt.Sprintf("%s", vi)
	}
	return stmt
}

// Where :
func (stmt *UpdateStmt) Where(fields ...any) *UpdateStmt {
	// stmt.Conditions = expr.And(fields...)
	return stmt
}

// Set :
func (stmt *UpdateStmt) Set(values ...primitive.KV) *UpdateStmt {
	stmt.Values = append(stmt.Values, values...)
	return stmt
}

// OrderBy :
func (stmt *UpdateStmt) OrderBy(fields ...any) *UpdateStmt {
	stmt.Sorts = fields
	return stmt
}

// Limit :
func (stmt *UpdateStmt) Limit(num uint) *UpdateStmt {
	if num > 0 {
		stmt.RowCount = num
	}
	return stmt
}
