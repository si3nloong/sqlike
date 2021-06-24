package casbin

import (
	"github.com/si3nloong/sqlike/v2/sql/expr"
	"github.com/si3nloong/sqlike/v2/x/primitive"
)

// Filter :
func Filter(fields ...interface{}) primitive.Group {
	return expr.And(fields...)
}
