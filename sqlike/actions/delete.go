package actions

import (
	"bitbucket.org/SianLoong/sqlike/sqlike/primitive"
	"bitbucket.org/SianLoong/sqlike/sqlike/sql/expr"
)

// DeleteStatement :
type DeleteStatement interface {
	Where(fields ...interface{}) DeleteStatement
	OrderBy(fields ...primitive.Sort) DeleteStatement
	Limit(num int) DeleteStatement
}

// DeleteActions :
type DeleteActions struct {
	Table      string
	Conditions []interface{}
	Sorts      []primitive.Sort
	Record     int
}

// Where :
func (f *DeleteActions) Where(fields ...interface{}) DeleteStatement {
	f.Conditions = expr.And(fields...)
	return f
}

// OrderBy :
func (f *DeleteActions) OrderBy(fields ...primitive.Sort) DeleteStatement {
	f.Sorts = fields
	return f
}

// Limit :
func (f *DeleteActions) Limit(num int) DeleteStatement {
	if num > 0 {
		f.Record = num
	}
	return f
}
