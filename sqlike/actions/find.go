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
	OrderBy(fields ...primitive.Sort) SelectStatement
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
	Sorts       []primitive.Sort
	Skip        uint
	Record      uint
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
func (f *FindActions) From(values ...string) SelectStatement {
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
		f.Skip = num
	}
	return f
}
