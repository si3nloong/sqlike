package actions

import (
	"strings"

	"github.com/si3nloong/sqlike/sql/expr"
	"github.com/si3nloong/sqlike/sqlike/primitive"
)

// SelectOneStatement :
type SelectOneStatement interface {
	Distinct() SelectOneStatement
	Select(fields ...interface{}) SelectOneStatement
	From(values ...string) SelectOneStatement
	Where(fields ...interface{}) SelectOneStatement
	Having(fields ...interface{}) SelectOneStatement
	GroupBy(fields ...interface{}) SelectOneStatement
	OrderBy(fields ...primitive.Sort) SelectOneStatement
}

// FindOneActions :
type FindOneActions struct {
	FindActions
}

// Select :
func (f *FindOneActions) Select(fields ...interface{}) SelectOneStatement {
	f.Projections = fields
	return f
}

// Distinct :
func (f *FindOneActions) Distinct() SelectOneStatement {
	f.DistinctOn = true
	return f
}

// From :
func (f *FindOneActions) From(values ...string) SelectOneStatement {
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
func (f *FindOneActions) Where(fields ...interface{}) SelectOneStatement {
	f.Conditions = expr.And(fields...)
	return f
}

// Having :
func (f *FindOneActions) Having(fields ...interface{}) SelectOneStatement {
	f.Havings = expr.And(fields...)
	return f
}

// OrderBy :
func (f *FindOneActions) OrderBy(fields ...primitive.Sort) SelectOneStatement {
	f.Sorts = fields
	return f
}

// GroupBy :
func (f *FindOneActions) GroupBy(fields ...interface{}) SelectOneStatement {
	f.GroupBys = fields
	return f
}
