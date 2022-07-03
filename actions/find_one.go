package actions

import (
	"github.com/si3nloong/sqlike/v2/sql/expr"
)

// SelectOneStatement :
type SelectOneStatement interface {
	Distinct() SelectOneStatement
	Select(fields ...any) SelectOneStatement
	Where(fields ...any) SelectOneStatement
	Having(fields ...any) SelectOneStatement
	GroupBy(fields ...any) SelectOneStatement
	OrderBy(fields ...any) SelectOneStatement
}

// FindOneActions :
type FindOneActions struct {
	FindActions
}

// Select :
func (act *FindOneActions) Select(fields ...any) SelectOneStatement {
	act.Projections = fields
	return act
}

// Distinct :
func (act *FindOneActions) Distinct() SelectOneStatement {
	act.DistinctOn = true
	return act
}

// Where :
func (act *FindOneActions) Where(fields ...any) SelectOneStatement {
	act.Conditions = expr.And(fields...)
	return act
}

// Having :
func (act *FindOneActions) Having(fields ...any) SelectOneStatement {
	act.Havings = expr.And(fields...)
	return act
}

// OrderBy :
func (act *FindOneActions) OrderBy(fields ...any) SelectOneStatement {
	act.Sorts = fields
	return act
}

// GroupBy :
func (act *FindOneActions) GroupBy(fields ...any) SelectOneStatement {
	act.GroupBys = fields
	return act
}
