package actions

import (
	"strings"

	"github.com/si3nloong/sqlike/sql/expr"
	"github.com/si3nloong/sqlike/sqlike/primitive"
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
	DistinctOn  bool
	Database    string
	Table       string
	Projections []interface{}
	IndexHints  string
	Conditions  primitive.Group
	Havings     primitive.Group
	GroupBys    []interface{}
	Sorts       []interface{}
	Skip        uint
	Count       uint
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
	switch length {
	case 1:
		act.Table = strings.TrimSpace(values[0])
	case 2:
		act.Database = strings.TrimSpace(values[0])
		act.Table = strings.TrimSpace(values[1])
	case 3:
		act.Database = strings.TrimSpace(values[0])
		act.Table = strings.TrimSpace(values[1])
		act.IndexHints = strings.TrimSpace(values[3])
	default:
		panic("invalid length of arguments")
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
		act.Count = num
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
