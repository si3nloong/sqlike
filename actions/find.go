package actions

import (
	"github.com/si3nloong/sqlike/v2/sql/expr"
	"github.com/si3nloong/sqlike/v2/x/primitive"
)

// SelectStatement :
type SelectStatement interface {
	Distinct() SelectStatement
	Select(fields ...any) SelectStatement
	Where(fields ...any) SelectStatement
	Having(fields ...any) SelectStatement
	GroupBy(fields ...any) SelectStatement
	OrderBy(fields ...any) SelectStatement
	Limit(num uint) SelectStatement
	Offset(num uint) SelectStatement
}

// FindActions :
type FindActions struct {
	DistinctOn  bool
	Database    string
	Table       string
	Projections []any
	IndexHints  string
	Conditions  primitive.Group
	Havings     primitive.Group
	GroupBys    []any
	Sorts       []any
	Skip        uint
	RowCount    uint
}

// Select :
func (act *FindActions) Select(fields ...any) SelectStatement {
	act.Projections = fields
	return act
}

// Distinct :
func (act *FindActions) Distinct() SelectStatement {
	act.DistinctOn = true
	return act
}

// Where :
func (act *FindActions) Where(fields ...any) SelectStatement {
	act.Conditions = expr.And(fields...)
	return act
}

// Having :
func (act *FindActions) Having(fields ...any) SelectStatement {
	act.Havings = expr.And(fields...)
	return act
}

// OrderBy :
func (act *FindActions) OrderBy(fields ...any) SelectStatement {
	act.Sorts = fields
	return act
}

// GroupBy :
func (act *FindActions) GroupBy(fields ...any) SelectStatement {
	act.GroupBys = fields
	return act
}

// Limit :
func (act *FindActions) Limit(num uint) SelectStatement {
	if num > 0 {
		act.RowCount = num
	}
	return act
}

// Offset :
func (act *FindActions) Offset(num uint) SelectStatement {
	if num > 0 {
		act.Skip = num
	}
	return act
}
