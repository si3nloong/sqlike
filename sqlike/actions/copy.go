package actions

import (
	"strings"

	"github.com/si3nloong/sqlike/sqlike/primitive"
	"github.com/si3nloong/sqlike/sql/expr"
)

// CopyStatement :
type CopyStatement interface {
	Distinct() CopyStatement
	Select(fields ...interface{}) CopyStatement
	From(table string) CopyStatement
	Where(fields ...interface{}) CopyStatement
	Having(fields ...interface{}) CopyStatement
	GroupBy(fields ...interface{}) CopyStatement
	OrderBy(fields ...primitive.Sort) CopyStatement
	Limit(num uint) CopyStatement
	Offset(num uint) CopyStatement
}

// CopyActions :
type CopyActions struct {
	FindActions
}

// Select :
func (f *CopyActions) Select(fields ...interface{}) CopyStatement {
	f.Projections = fields
	return f
}

// Distinct :
func (f *CopyActions) Distinct() CopyStatement {
	f.DistinctOn = true
	return f
}

// From :
func (f *CopyActions) From(table string) CopyStatement {
	table = strings.TrimSpace(table)
	if table != "" {
		f.Table = table
	}
	return f
}

// Where :
func (f *CopyActions) Where(fields ...interface{}) CopyStatement {
	f.Conditions = expr.And(fields...)
	return f
}

// Having :
func (f *CopyActions) Having(fields ...interface{}) CopyStatement {
	f.Havings = expr.And(fields...)
	return f
}

// OrderBy :
func (f *CopyActions) OrderBy(fields ...primitive.Sort) CopyStatement {
	f.Sorts = fields
	return f
}

// GroupBy :
func (f *CopyActions) GroupBy(fields ...interface{}) CopyStatement {
	f.GroupBys = fields
	return f
}

// Limit :
func (f *CopyActions) Limit(num uint) CopyStatement {
	if num > 0 {
		f.Record = num
	}
	return f
}

// Offset :
func (f *CopyActions) Offset(num uint) CopyStatement {
	if num > 0 {
		f.Skip = num
	}
	return f
}
