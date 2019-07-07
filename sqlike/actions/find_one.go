package actions

import (
	"strings"

	"github.com/si3nloong/sqlike/sqlike/primitive"
	"github.com/si3nloong/sqlike/sql/expr"
)

// SelectOneStatement :
type SelectOneStatement interface {
	Distinct() SelectOneStatement
	Select(fields ...interface{}) SelectOneStatement
	From(table string) SelectOneStatement
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
func (f *FindOneActions) From(table string) SelectOneStatement {
	table = strings.TrimSpace(table)
	if table != "" {
		f.Table = table
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
