package actions

import (
	"strings"

	"github.com/si3nloong/sqlike/sqlike/primitive"
	"github.com/si3nloong/sqlike/sqlike/sql/expr"
)

// SelectStatement :
type SelectStatement interface {
	Distinct() SelectStatement
	Select(fields ...interface{}) SelectStatement
	From(table string) SelectStatement
	Where(fields ...interface{}) SelectStatement
	Having(fields ...interface{}) SelectStatement
	GroupBy(fields ...interface{}) SelectStatement
	OrderBy(fields ...primitive.Sort) SelectStatement
	Limit(num uint) SelectStatement
	Offset(num uint) SelectStatement
}

// FindActions :
type FindActions struct {
	Table       string
	DistinctOn  bool
	Projections []interface{}
	Conditions  []interface{}
	Havings     []interface{}
	GroupBys    []interface{}
	Sorts       []primitive.Sort
	Record      uint
	Offs        uint
}

// Select :
func (f *FindActions) Select(fields ...interface{}) SelectStatement {
	f.Projections = fields
	return f
}

// Distinct :
func (f *FindActions) Distinct() SelectStatement {
	f.DistinctOn = true
	return f
}

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
func (f *FindActions) Limit(num uint) SelectStatement {
	if num > 0 {
		f.Record = num
	}
	return f
}

// Offset :
func (f *FindActions) Offset(num uint) SelectStatement {
	if num > 0 {
		f.Offs = num
	}
	return f
}
