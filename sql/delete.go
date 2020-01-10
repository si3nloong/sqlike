package sql

import (
	"github.com/si3nloong/sqlike/sqlike/primitive"
)

// DeleteStmt :
type DeleteStmt struct {
	Tables     []interface{}
	Conditions primitive.Group
	Sorts      []interface{}
	Max        uint
}

// From :
func (stmt *DeleteStmt) From() *DeleteStmt {
	return stmt
}

// Where :
func (stmt *DeleteStmt) Where(fields ...interface{}) *DeleteStmt {
	// stmt.Conditions = expr.And(fields...)
	return stmt
}

// OrderBy :
func (stmt *DeleteStmt) OrderBy(fields ...interface{}) *DeleteStmt {
	stmt.Sorts = fields
	return stmt
}

// Limit :
func (stmt *DeleteStmt) Limit(num uint) *DeleteStmt {
	if num > 0 {
		stmt.Max = num
	}
	return stmt
}
