package actions

import (
	"github.com/si3nloong/sqlike/sql/expr"
	"github.com/si3nloong/sqlike/sqlike/primitive"
)

// UpdateOneStatement :
type UpdateOneStatement interface {
	Where(fields ...interface{}) UpdateOneStatement
	Set(values ...primitive.KV) UpdateOneStatement
	OrderBy(fields ...interface{}) UpdateOneStatement
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
func (f *UpdateOneActions) Set(values ...primitive.KV) UpdateOneStatement {
	f.Values = append(f.Values, values...)
	return f
}

// OrderBy :
func (f *UpdateOneActions) OrderBy(fields ...interface{}) UpdateOneStatement {
	f.Sorts = fields
	return f
}
