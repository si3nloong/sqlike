package actions

import (
	"github.com/si3nloong/sqlike/v2/sql/expr"
	"github.com/si3nloong/sqlike/v2/x/primitive"
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
func (act *UpdateOneActions) Where(fields ...interface{}) UpdateOneStatement {
	act.Conditions = expr.And(fields...).Values
	return act
}

// Set :
func (act *UpdateOneActions) Set(values ...primitive.KV) UpdateOneStatement {
	act.Values = append(act.Values, values...)
	return act
}

// OrderBy :
func (act *UpdateOneActions) OrderBy(fields ...interface{}) UpdateOneStatement {
	act.Sorts = fields
	return act
}
