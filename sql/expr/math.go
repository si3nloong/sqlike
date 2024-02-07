package expr

import (
	"golang.org/x/exp/constraints"

	"github.com/si3nloong/sqlike/v2/internal/primitive"
)

// Increment :
func Increment[V constraints.Unsigned](field string, value V) primitive.Math {
	return primitive.Math{
		Field: field,
		Mode:  primitive.Add,
		Value: uint(value),
	}
}

// Decrement :
func Decrement[V constraints.Unsigned](field string, value V) primitive.Math {
	return primitive.Math{
		Field: field,
		Mode:  primitive.Deduct,
		Value: uint(value),
	}
}

// Multiply :
func Multiply(fields ...any) (grp primitive.Group) {
	for i, f := range fields {
		if i > 0 {
			grp.Values = append(grp.Values, Raw("*"))
		}
		grp.Values = append(grp.Values, wrapColumn(f))
	}
	return
}
