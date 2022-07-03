package actions

import (
	"github.com/si3nloong/sqlike/v2/sql/expr"
	"github.com/si3nloong/sqlike/v2/x/primitive"
)

// UpdateOneStatement :
type UpdateOneStatement interface {
	Where(fields ...any) UpdateOneStatement
	Set(values ...primitive.KV) UpdateOneStatement
	OrderBy(fields ...any) UpdateOneStatement
}

// UpdateOneActions :
type UpdateOneActions struct {
	UpdateActions
}

// Where :
func (act *UpdateOneActions) Where(fields ...any) UpdateOneStatement {
	act.Conditions = expr.And(fields...).Values
	return act
}

// Set :
func (act *UpdateOneActions) Set(values ...primitive.KV) UpdateOneStatement {
	act.Values = append(act.Values, values...)
	return act
}

// OrderBy :
func (act *UpdateOneActions) OrderBy(fields ...any) UpdateOneStatement {
	act.Sorts = fields
	return act
}
