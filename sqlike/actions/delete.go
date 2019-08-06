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
func (f *DeleteActions) Where(fields ...interface{}) DeleteStatement {
	f.Conditions = expr.And(fields...)
	return f
}

// OrderBy :
func (f *DeleteActions) OrderBy(fields ...interface{}) DeleteStatement {
	f.Sorts = fields
	return f
}

// Limit :
func (f *DeleteActions) Limit(num uint) DeleteStatement {
	if num > 0 {
		f.Record = num
	}
	return f
}
