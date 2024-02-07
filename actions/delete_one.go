package actions

import (
	"github.com/si3nloong/sqlike/v2/sql/expr"
)

// DeleteOneStatement :
type DeleteOneStatement interface {
	Where(fields ...any) DeleteOneStatement
	OrderBy(fields ...any) DeleteOneStatement
}

// DeleteOneActions :
type DeleteOneActions struct {
	DeleteActions
}

// Where :
func (act *DeleteOneActions) Where(fields ...any) DeleteOneStatement {
	act.Conditions = expr.And(fields...).Values
	return act
}

// OrderBy :
func (act *DeleteOneActions) OrderBy(fields ...any) DeleteOneStatement {
	act.Sorts = fields
	return act
}

// Limit :
func (act *DeleteOneActions) Limit(num uint) DeleteOneStatement {
	if num > 0 {
		act.RowCount = num
	}
	return act
}
