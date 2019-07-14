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
	From(table string) PaginateStatement
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
func (f *PaginateActions) From(table string) PaginateStatement {
	table = strings.TrimSpace(table)
	if table != "" {
		f.Table = table
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
