package actions

import (
	"github.com/si3nloong/sqlike/v2/sql/expr"
	"github.com/si3nloong/sqlike/v2/x/primitive"
)

// UpdateStatement :
type UpdateStatement interface {
	Where(fields ...interface{}) UpdateStatement
	Set(values ...primitive.KV) UpdateStatement
	OrderBy(fields ...interface{}) UpdateStatement
	Limit(num uint) UpdateStatement
}

// UpdateActions :
type UpdateActions struct {
	Database   string
	Table      string
	Conditions []interface{}
	Values     []primitive.KV
	Sorts      []interface{}
	RowCount   uint
}

// Where :
func (act *UpdateActions) Where(fields ...interface{}) UpdateStatement {
	act.Conditions = expr.And(fields...).Values
	return act
}

// Set :
func (act *UpdateActions) Set(values ...primitive.KV) UpdateStatement {
	act.Values = append(act.Values, values...)
	return act
}

// OrderBy :
func (act *UpdateActions) OrderBy(fields ...interface{}) UpdateStatement {
	act.Sorts = fields
	return act
}

// Limit :
func (act *UpdateActions) Limit(num uint) UpdateStatement {
	if num > 0 {
		act.RowCount = num
	}
	return act
}
