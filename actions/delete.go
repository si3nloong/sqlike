package actions

import (
	"github.com/si3nloong/sqlike/v2/sql/expr"
)

// DeleteStatement :
type DeleteStatement interface {
	Where(fields ...interface{}) DeleteStatement
	OrderBy(fields ...interface{}) DeleteStatement
	Limit(num uint) DeleteStatement
}

// DeleteActions :
type DeleteActions struct {
	Database   string
	Table      string
	Conditions []interface{}
	Sorts      []interface{}
	RowCount   uint
}

// Where :
func (act *DeleteActions) Where(fields ...interface{}) DeleteStatement {
	act.Conditions = expr.And(fields...).Values
	return act
}

// OrderBy :
func (act *DeleteActions) OrderBy(fields ...interface{}) DeleteStatement {
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