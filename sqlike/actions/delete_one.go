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
func (f *DeleteOneActions) Where(fields ...interface{}) DeleteOneStatement {
	f.Conditions = expr.And(fields...)
	return f
}

// OrderBy :
func (f *DeleteOneActions) OrderBy(fields ...interface{}) DeleteOneStatement {
	f.Sorts = fields
	return f
}

// Limit :
func (f *DeleteOneActions) Limit(num uint) DeleteOneStatement {
	if num > 0 {
		f.Record = num
	}
	return f
}
