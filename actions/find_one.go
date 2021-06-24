package actions

import (
	"strings"

	"github.com/si3nloong/sqlike/v2/sql/expr"
)

// SelectOneStatement :
type SelectOneStatement interface {
	Distinct() SelectOneStatement
	Select(fields ...interface{}) SelectOneStatement
	From(values ...string) SelectOneStatement
	Where(fields ...interface{}) SelectOneStatement
	Having(fields ...interface{}) SelectOneStatement
	GroupBy(fields ...interface{}) SelectOneStatement
	OrderBy(fields ...interface{}) SelectOneStatement
}

// FindOneActions :
type FindOneActions struct {
	FindActions
}

// Select :
func (act *FindOneActions) Select(fields ...interface{}) SelectOneStatement {
	act.Projections = fields
	return act
}

// Distinct :
func (act *FindOneActions) Distinct() SelectOneStatement {
	act.DistinctOn = true
	return act
}

// From :
func (act *FindOneActions) From(values ...string) SelectOneStatement {
	length := len(values)
	if length == 0 {
		panic("empty table name")
	}
	if length > 0 {
		act.Table = strings.TrimSpace(values[0])
	}
	if length > 1 {
		act.Database = strings.TrimSpace(values[0])
		act.Table = strings.TrimSpace(values[1])
	}
	return act
}

// Where :
func (act *FindOneActions) Where(fields ...interface{}) SelectOneStatement {
	act.Conditions = expr.And(fields...)
	return act
}

// Having :
func (act *FindOneActions) Having(fields ...interface{}) SelectOneStatement {
	act.Havings = expr.And(fields...)
	return act
}

// OrderBy :
func (act *FindOneActions) OrderBy(fields ...interface{}) SelectOneStatement {
	act.Sorts = fields
	return act
}

// GroupBy :
func (act *FindOneActions) GroupBy(fields ...interface{}) SelectOneStatement {
	act.GroupBys = fields
	return act
}
