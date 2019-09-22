package actions

import (
	"github.com/si3nloong/sqlike/sql/expr"
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
	Record     uint
}

// Where :
func (act *DeleteActions) Where(fields ...interface{}) DeleteStatement {
	act.Conditions = expr.And(fields...)
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
		act.Record = num
	}
	return act
}
