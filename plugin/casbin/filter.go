package casbin

import (
	"github.com/si3nloong/sqlike/sql/expr"
	"github.com/si3nloong/sqlike/sqlike/primitive"
)

// Filter :
func Filter(fields ...interface{}) primitive.Group {
	return expr.And(fields...)
}
