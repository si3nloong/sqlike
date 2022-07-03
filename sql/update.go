package sql

import (
	"github.com/si3nloong/sqlike/v2/x/primitive"
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
func Update(tables ...any) *UpdateStmt {
	stmt := new(UpdateStmt)
	return stmt.Update(tables...)
}

// Update :
func (stmt *UpdateStmt) Update(fields ...any) *UpdateStmt {
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
