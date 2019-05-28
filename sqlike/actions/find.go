package actions

import (
	"strings"

	"github.com/si3nloong/sqlike/sqlike/primitive"
	"github.com/si3nloong/sqlike/sqlike/sql/expr"
)

// SelectStatement :
type SelectStatement interface {
	Select(fields ...interface{}) SelectStatement
	From(table string) SelectStatement
	Where(fields ...interface{}) SelectStatement
	Having(fields ...interface{}) SelectStatement
	GroupBy(fields ...interface{}) SelectStatement
	OrderBy(fields ...primitive.Sort) SelectStatement
	Limit(num int) SelectStatement
	Offset(num int) SelectStatement
}

// FindActions :
type FindActions struct {
	Table       string
	Projections []interface{}
	Conditions  []interface{}
	Havings     []interface{}
	GroupBys    []interface{}
	Sorts       []primitive.Sort
	Record      int
	Offs        int
}

// Select :
func (f *FindActions) Select(fields ...interface{}) SelectStatement {
	f.Projections = fields
	return f
}

// Distinct :
// func (f *FindActions) Distinct() *FindActions {
// 	return f
// }

// From :
func (f *FindActions) From(table string) SelectStatement {
	table = strings.TrimSpace(table)
	if table != "" {
		f.Table = table
	}
	return f
}

// Where :
func (f *FindActions) Where(fields ...interface{}) SelectStatement {
	f.Conditions = expr.And(fields...)
	return f
}

// Having :
func (f *FindActions) Having(fields ...interface{}) SelectStatement {
	f.Havings = expr.And(fields...)
	return f
}

// OrderBy :
func (f *FindActions) OrderBy(fields ...primitive.Sort) SelectStatement {
	f.Sorts = fields
	return f
}

// GroupBy :
func (f *FindActions) GroupBy(fields ...interface{}) SelectStatement {
	f.GroupBys = fields
	return f
}

// Limit :
func (f *FindActions) Limit(num int) SelectStatement {
	if num > 0 {
		f.Record = num
	}
	return f
}

// Offset :
func (f *FindActions) Offset(num int) SelectStatement {
	if num > 0 {
		f.Offs = num
	}
	return f
}
