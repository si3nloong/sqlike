package sql

import (
	"reflect"

	"github.com/si3nloong/sqlike/v2/sql/expr"
	"github.com/si3nloong/sqlike/v2/x/primitive"
	"github.com/si3nloong/sqlike/v2/x/reflext"
)

// SelectStmt :
type SelectStmt struct {
	DistinctOn bool
	Tables     []any
	Exprs      []any
	Joins      []any
	IndexHints string
	Conditions primitive.Group
	Havings    primitive.Group
	Groups     []any
	Sorts      []any
	RowCount   uint
	Skip       uint
	Opts       []any
}

// Select :
func Select(fields ...any) *SelectStmt {
	stmt := new(SelectStmt)
	return stmt.Select(fields...)
}

// Select :
func (stmt *SelectStmt) Select(fields ...any) *SelectStmt {
	if len(fields) == 1 {
		switch fields[0].(type) {
		case primitive.As, *SelectStmt:
			grp := primitive.Group{}
			grp.Values = append(grp.Values, primitive.Raw{Value: "("})
			grp.Values = append(grp.Values, fields...)
			grp.Values = append(grp.Values, primitive.Raw{Value: ")"})
			stmt.Exprs = append(stmt.Exprs, grp)
		default:
			stmt.Exprs = append(stmt.Exprs, fields...)
		}
		return stmt
	}
	stmt.Exprs = fields
	return stmt
}

// From :
func (stmt *SelectStmt) From(values ...any) *SelectStmt {
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

// Where :
func (stmt *SelectStmt) Where(fields ...any) *SelectStmt {
	stmt.Conditions = expr.And(fields...)
	return stmt
}

// Having :
func (stmt *SelectStmt) Having(fields ...any) *SelectStmt {
	stmt.Havings = expr.And(fields...)
	return stmt
}

// OrderBy :
func (stmt *SelectStmt) OrderBy(fields ...any) *SelectStmt {
	stmt.Sorts = fields
	return stmt
}

// GroupBy :
func (stmt *SelectStmt) GroupBy(fields ...any) *SelectStmt {
	stmt.Groups = fields
	return stmt
}

// Limit :
func (stmt *SelectStmt) Limit(num uint) *SelectStmt {
	if num > 0 {
		stmt.RowCount = num
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

// Option :
func (stmt *SelectStmt) Option(opts ...any) *SelectStmt {
	stmt.Opts = opts
	return stmt
}

func mustString(it any) string {
	v := reflext.Indirect(reflext.ValueOf(it))
	if v.Kind() != reflect.String {
		panic("unsupported data type")
	}
	return v.String()
}
