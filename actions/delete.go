package actions

import (
	"github.com/si3nloong/sqlike/v2/sql/expr"
)

// DeleteStatement :
type DeleteStatement interface {
	Where(fields ...any) DeleteStatement
	OrderBy(fields ...any) DeleteStatement
	Limit(num uint) DeleteStatement
}

// DeleteActions :
type DeleteActions struct {
	Database   string
	Table      string
	Conditions []any
	Sorts      []any
	RowCount   uint
}

// Where :
func (act *DeleteActions) Where(fields ...any) DeleteStatement {
	act.Conditions = expr.And(fields...).Values
	return act
}

// OrderBy :
func (act *DeleteActions) OrderBy(fields ...any) DeleteStatement {
	act.Sorts = fields
	return act
}

// Limit :
func (act *DeleteActions) Limit(num uint) DeleteStatement {
	if num > 0 {
		act.RowCount = num
	}
	return act
}
