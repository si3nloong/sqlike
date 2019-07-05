package actions

import (
	"github.com/si3nloong/sqlike/sqlike/primitive"
	"github.com/si3nloong/sqlike/sqlike/sql/expr"
)

// UpdateStatement :
type UpdateStatement interface {
	Where(fields ...interface{}) UpdateStatement
	Set(values ...primitive.KV) UpdateStatement
	OrderBy(fields ...primitive.Sort) UpdateStatement
	Limit(num uint) UpdateStatement
}

// UpdateActions :
type UpdateActions struct {
	Table      string
	Conditions []interface{}
	Values     []primitive.KV
	Sorts      []primitive.Sort
	Record     uint
}

// Where :
func (f *UpdateActions) Where(fields ...interface{}) UpdateStatement {
	f.Conditions = expr.And(fields...)
	return f
}

// Set :
func (f *UpdateActions) Set(values ...primitive.KV) UpdateStatement {
	f.Values = append(f.Values, values...)
	return f
}

// OrderBy :
func (f *UpdateActions) OrderBy(fields ...primitive.Sort) UpdateStatement {
	f.Sorts = fields
	return f
}

// Limit :
func (f *UpdateActions) Limit(num uint) UpdateStatement {
	if num > 0 {
		f.Record = num
	}
	return f
}
