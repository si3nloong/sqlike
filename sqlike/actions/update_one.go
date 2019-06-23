package actions

import (
	"github.com/si3nloong/sqlike/sqlike/primitive"
	"github.com/si3nloong/sqlike/sqlike/sql/expr"
)

// UpdateOneStatement :
type UpdateOneStatement interface {
	Where(fields ...interface{}) UpdateOneStatement
	Set(field string, value interface{}) UpdateOneStatement
	OrderBy(fields ...primitive.Sort) UpdateOneStatement
}

// UpdateOneActions :
type UpdateOneActions struct {
	UpdateActions
}

// Where :
func (f *UpdateOneActions) Where(fields ...interface{}) UpdateOneStatement {
	f.Conditions = expr.And(fields...)
	return f
}

// Set :
func (f *UpdateOneActions) Set(field string, value interface{}) UpdateOneStatement {
	f.Values = append(f.Values, primitive.C{
		Field:    primitive.Col(field),
		Operator: primitive.Equal,
		Value:    value,
	})
	return f
}

// OrderBy :
func (f *UpdateOneActions) OrderBy(fields ...primitive.Sort) UpdateOneStatement {
	f.Sorts = fields
	return f
}
