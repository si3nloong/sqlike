package actions

import (
	"github.com/si3nloong/sqlike/sql/expr"
)

// DeleteOneStatement :
type DeleteOneStatement interface {
	Where(fields ...interface{}) DeleteOneStatement
	OrderBy(fields ...interface{}) DeleteOneStatement
}

// DeleteOneActions :
type DeleteOneActions struct {
	DeleteActions
}

// Where :
func (act *DeleteOneActions) Where(fields ...interface{}) DeleteOneStatement {
	act.Conditions = expr.And(fields...).Values
	return act
}

// OrderBy :
func (act *DeleteOneActions) OrderBy(fields ...interface{}) DeleteOneStatement {
	act.Sorts = fields
	return act
}

// Limit :
func (act *DeleteOneActions) Limit(num uint) DeleteOneStatement {
	if num > 0 {
		act.Record = num
	}
	return act
}
