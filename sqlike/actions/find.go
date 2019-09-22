package actions

import (
	"strings"

	"github.com/si3nloong/sqlike/sql/expr"
)

// SelectStatement :
type SelectStatement interface {
	Distinct() SelectStatement
	Select(fields ...interface{}) SelectStatement
	From(values ...string) SelectStatement
	Where(fields ...interface{}) SelectStatement
	Having(fields ...interface{}) SelectStatement
	GroupBy(fields ...interface{}) SelectStatement
	OrderBy(fields ...interface{}) SelectStatement
	Limit(num uint) SelectStatement
	Offset(num uint) SelectStatement
}

// FindActions :
type FindActions struct {
	Database    string
	Table       string
	DistinctOn  bool
	Projections []interface{}
	Conditions  []interface{}
	Havings     []interface{}
	GroupBys    []interface{}
	Sorts       []interface{}
	Skip        uint
	Record      uint
}

// Select :
func (act *FindActions) Select(fields ...interface{}) SelectStatement {
	act.Projections = fields
	return act
}

// Distinct :
func (act *FindActions) Distinct() SelectStatement {
	act.DistinctOn = true
	return act
}

// From :
func (act *FindActions) From(values ...string) SelectStatement {
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
func (act *FindActions) Where(fields ...interface{}) SelectStatement {
	act.Conditions = expr.And(fields...)
	return act
}

// Having :
func (act *FindActions) Having(fields ...interface{}) SelectStatement {
	act.Havings = expr.And(fields...)
	return act
}

// OrderBy :
func (act *FindActions) OrderBy(fields ...interface{}) SelectStatement {
	act.Sorts = fields
	return act
}

// GroupBy :
func (act *FindActions) GroupBy(fields ...interface{}) SelectStatement {
	act.GroupBys = fields
	return act
}

// Limit :
func (act *FindActions) Limit(num uint) SelectStatement {
	if num > 0 {
		act.Record = num
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
