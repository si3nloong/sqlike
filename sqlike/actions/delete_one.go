package actions

import (
	"github.com/si3nloong/sqlike/sql/expr"
	"github.com/si3nloong/sqlike/sqlike/primitive"
)

// DeleteOneStatement :
type DeleteOneStatement interface {
	Where(fields ...interface{}) DeleteOneStatement
	OrderBy(fields ...primitive.Sort) DeleteOneStatement
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
func (f *DeleteOneActions) OrderBy(fields ...primitive.Sort) DeleteOneStatement {
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
