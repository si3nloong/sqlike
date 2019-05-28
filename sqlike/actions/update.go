package actions

import (
	"github.com/si3nloong/sqlike/sqlike/primitive"
	"github.com/si3nloong/sqlike/sqlike/sql/expr"
)

// UpdateStatement :
type UpdateStatement interface {
	Where(fields ...interface{}) UpdateStatement
	Set(field string, value interface{}) UpdateStatement
	OrderBy(fields ...primitive.Sort) UpdateStatement
	Limit(num int) UpdateStatement
}

// UpdateActions :
type UpdateActions struct {
	Table      string
	Conditions []interface{}
	Values     []primitive.C
	Sorts      []primitive.Sort
	Record     int
}

// Where :
func (f *UpdateActions) Where(fields ...interface{}) UpdateStatement {
	f.Conditions = expr.And(fields...)
	return f
}

// Set :
func (f *UpdateActions) Set(field string, value interface{}) UpdateStatement {
	f.Values = append(f.Values, primitive.C{
		Field:    primitive.Col(field),
		Operator: primitive.Equal,
		Values:   []interface{}{value},
	})
	return f
}

// OrderBy :
func (f *UpdateActions) OrderBy(fields ...primitive.Sort) UpdateStatement {
	f.Sorts = fields
	return f
}

// Limit :
func (f *UpdateActions) Limit(num int) UpdateStatement {
	if num > 0 {
		f.Record = num
	}
	return f
}
