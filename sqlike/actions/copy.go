package actions

import (
	"strings"

	"github.com/si3nloong/sqlike/sql/expr"
)

// CopyStatement :
type CopyStatement interface {
	Distinct() CopyStatement
	Select(fields ...interface{}) CopyStatement
	From(values ...string) CopyStatement
	Where(fields ...interface{}) CopyStatement
	Having(fields ...interface{}) CopyStatement
	GroupBy(fields ...interface{}) CopyStatement
	OrderBy(fields ...interface{}) CopyStatement
	Limit(num uint) CopyStatement
	Offset(num uint) CopyStatement
}

// CopyActions :
type CopyActions struct {
	FindActions
}

// Select :
func (act *CopyActions) Select(fields ...interface{}) CopyStatement {
	act.Projections = fields
	return act
}

// Distinct :
func (act *CopyActions) Distinct() CopyStatement {
	act.DistinctOn = true
	return act
}

// From :
func (act *CopyActions) From(values ...string) CopyStatement {
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
func (act *CopyActions) Where(fields ...interface{}) CopyStatement {
	act.Conditions = expr.And(fields...)
	return act
}

// Having :
func (act *CopyActions) Having(fields ...interface{}) CopyStatement {
	act.Havings = expr.And(fields...)
	return act
}

// OrderBy :
func (act *CopyActions) OrderBy(fields ...interface{}) CopyStatement {
	act.Sorts = fields
	return act
}

// GroupBy :
func (act *CopyActions) GroupBy(fields ...interface{}) CopyStatement {
	act.GroupBys = fields
	return act
}

// Limit :
func (act *CopyActions) Limit(num uint) CopyStatement {
	if num > 0 {
		act.Count = num
	}
	return act
}

// Offset :
func (act *CopyActions) Offset(num uint) CopyStatement {
	if num > 0 {
		act.Skip = num
	}
	return act
}
