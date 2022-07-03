package sql

import (
	"github.com/si3nloong/sqlike/v2/sql/expr"
	"github.com/si3nloong/sqlike/v2/x/primitive"
)

// DeleteStmt :
type DeleteStmt struct {
	Tables     []any
	Conditions primitive.Group
	Sorts      []any
	RowCount   uint
}

// From :
func (stmt *DeleteStmt) From(values ...any) *DeleteStmt {
	length := len(values)
	if length == 0 {
		panic("empty table name")
	}
	switch length {
	case 1:
		stmt.Tables = append(stmt.Tables, values[0])
	case 2:
		stmt.Tables = append(stmt.Tables,
			primitive.Column{
				Table: mustString(values[0]),
				Name:  mustString(values[1]),
			},
		)
	default:
		panic("invalid length of arguments")
	}
	return stmt
}

// Where :
func (stmt *DeleteStmt) Where(fields ...any) *DeleteStmt {
	stmt.Conditions = expr.And(fields...)
	return stmt
}

// OrderBy :
func (stmt *DeleteStmt) OrderBy(fields ...any) *DeleteStmt {
	stmt.Sorts = fields
	return stmt
}

// Limit :
func (stmt *DeleteStmt) Limit(num uint) *DeleteStmt {
	if num > 0 {
		stmt.RowCount = num
	}
	return stmt
}
