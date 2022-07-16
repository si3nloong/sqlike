package casbin

import (
	"github.com/si3nloong/sqlike/v2/internal/primitive"
	"github.com/si3nloong/sqlike/v2/sql/expr"
)

// Filter :
func Filter(fields ...any) primitive.Group {
	return expr.And(fields...)
}
