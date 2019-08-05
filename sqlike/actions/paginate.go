package actions

import (
	"strings"

	"github.com/si3nloong/sqlike/sql/expr"
	"github.com/si3nloong/sqlike/sqlike/primitive"
)

// PaginateStatement :
type PaginateStatement interface {
	Distinct() PaginateStatement
	Select(fields ...interface{}) PaginateStatement
	From(values ...string) PaginateStatement
	Where(fields ...interface{}) PaginateStatement
	Having(fields ...interface{}) PaginateStatement
	GroupBy(fields ...interface{}) PaginateStatement
	OrderBy(fields ...primitive.Sort) PaginateStatement
	Limit(num uint) PaginateStatement
	Offset(num uint) PaginateStatement
}

// PaginateActions :
type PaginateActions struct {
	FindActions
}

// Select :
func (f *PaginateActions) Select(fields ...interface{}) PaginateStatement {
	f.Projections = fields
	return f
}

// Distinct :
func (f *PaginateActions) Distinct() PaginateStatement {
	f.DistinctOn = true
	return f
}

// From :
func (f *PaginateActions) From(values ...string) PaginateStatement {
	length := len(values)
	if length == 0 {
		panic("empty table name")
	}
	if length > 0 {
		f.Table = strings.TrimSpace(values[0])
	}
	if length > 1 {
		f.Database = strings.TrimSpace(values[0])
		f.Table = strings.TrimSpace(values[1])
	}
	return f
}

// Where :
func (f *PaginateActions) Where(fields ...interface{}) PaginateStatement {
	f.Conditions = expr.And(fields...)
	return f
}

// Having :
func (f *PaginateActions) Having(fields ...interface{}) PaginateStatement {
	f.Havings = expr.And(fields...)
	return f
}

// OrderBy :
func (f *PaginateActions) OrderBy(fields ...primitive.Sort) PaginateStatement {
	f.Sorts = fields
	return f
}

// GroupBy :
func (f *PaginateActions) GroupBy(fields ...interface{}) PaginateStatement {
	f.GroupBys = fields
	return f
}

// Limit :
func (f *PaginateActions) Limit(num uint) PaginateStatement {
	if num > 0 {
		f.Record = num
	}
	return f
}

// Offset :
func (f *PaginateActions) Offset(num uint) PaginateStatement {
	if num > 0 {
		f.Skip = num
	}
	return f
}
