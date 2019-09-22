package actions

import (
	"strings"

	"github.com/si3nloong/sqlike/sql/expr"
)

// PaginateStatement :
type PaginateStatement interface {
	Distinct() PaginateStatement
	Select(fields ...interface{}) PaginateStatement
	From(values ...string) PaginateStatement
	Where(fields ...interface{}) PaginateStatement
	Having(fields ...interface{}) PaginateStatement
	GroupBy(fields ...interface{}) PaginateStatement
	OrderBy(fields ...interface{}) PaginateStatement
	Limit(num uint) PaginateStatement
	Offset(num uint) PaginateStatement
}

// PaginateActions :
type PaginateActions struct {
	FindActions
}

// Select :
func (act *PaginateActions) Select(fields ...interface{}) PaginateStatement {
	act.Projections = fields
	return act
}

// Distinct :
func (act *PaginateActions) Distinct() PaginateStatement {
	act.DistinctOn = true
	return act
}

// From :
func (act *PaginateActions) From(values ...string) PaginateStatement {
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
func (act *PaginateActions) Where(fields ...interface{}) PaginateStatement {
	act.Conditions = expr.And(fields...)
	return act
}

// Having :
func (act *PaginateActions) Having(fields ...interface{}) PaginateStatement {
	act.Havings = expr.And(fields...)
	return act
}

// OrderBy :
func (act *PaginateActions) OrderBy(fields ...interface{}) PaginateStatement {
	act.Sorts = fields
	return act
}

// GroupBy :
func (act *PaginateActions) GroupBy(fields ...interface{}) PaginateStatement {
	act.GroupBys = fields
	return act
}

// Limit :
func (act *PaginateActions) Limit(num uint) PaginateStatement {
	if num > 0 {
		act.Record = num
	}
	return act
}

// Offset :
func (act *PaginateActions) Offset(num uint) PaginateStatement {
	if num > 0 {
		act.Skip = num
	}
	return act
}
