package sql

import (
	"reflect"

	"github.com/si3nloong/sqlike/reflext"
	"github.com/si3nloong/sqlike/sql/expr"
	"github.com/si3nloong/sqlike/sqlike/primitive"
)

// SelectStmt :
type SelectStmt struct {
	DistinctOn  bool
	Tables      []interface{}
	Projections []interface{}
	IndexHints  string
	Conditions  primitive.Group
	Havings     primitive.Group
	Groups      []interface{}
	Sorts       []interface{}
	Max         uint
	Skip        uint
}

// Distinct :
func (stmt *SelectStmt) Select(fields ...interface{}) *SelectStmt {
	if len(fields) == 1 {
		switch fields[0].(type) {
		case primitive.As, *SelectStmt:
			grp := primitive.Group{}
			grp.Values = append(grp.Values, expr.Raw("("))
			grp.Values = append(grp.Values, fields...)
			grp.Values = append(grp.Values, expr.Raw(")"))
			stmt.Projections = append(stmt.Projections, grp)
		default:
			stmt.Projections = append(stmt.Projections, fields...)
		}
		return stmt
	}
	stmt.Projections = fields
	return stmt
}

// From :
func (stmt *SelectStmt) From(values ...interface{}) *SelectStmt {
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
	case 3:
		stmt.Tables = append(stmt.Tables,
			primitive.Column{
				Table: mustString(values[0]),
				Name:  mustString(values[1]),
			},
			values[2],
		)
	default:
		panic("invalid length of arguments")
	}
	return stmt
}

// Distinct :
func (stmt *SelectStmt) Distinct() *SelectStmt {
	stmt.DistinctOn = true
	return stmt
}

func (stmt *SelectStmt) Where(fields ...interface{}) *SelectStmt {
	stmt.Conditions = expr.And(fields...)
	return stmt
}

func (stmt *SelectStmt) Having(fields ...interface{}) *SelectStmt {
	stmt.Havings = expr.And(fields...)
	return stmt
}

func (stmt *SelectStmt) OrderBy(fields ...interface{}) *SelectStmt {
	stmt.Sorts = fields
	return stmt
}

// GroupBy :
func (stmt *SelectStmt) GroupBy(fields ...interface{}) *SelectStmt {
	stmt.Groups = fields
	return stmt
}

// Limit :
func (stmt *SelectStmt) Limit(num uint) *SelectStmt {
	if num > 0 {
		stmt.Max = num
	}
	return stmt
}

// Offset :
func (stmt *SelectStmt) Offset(num uint) *SelectStmt {
	if num > 0 {
		stmt.Skip = num
	}
	return stmt
}

// Select :
func Select(fields ...interface{}) *SelectStmt {
	stmt := new(SelectStmt)
	return stmt.Select(fields...)
}

func mustString(it interface{}) string {
	v := reflext.Indirect(reflext.ValueOf(it))
	if v.Kind() != reflect.String {
		panic("unsupported data type")
	}
	return v.String()
}
