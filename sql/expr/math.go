package expr

import (
	"golang.org/x/exp/constraints"

	"github.com/si3nloong/sqlike/v2/x/primitive"
)

// Increment :
func Increment[T constraints.Integer](field string, value T) primitive.Math {
	return primitive.Math{
		Field: field,
		Mode:  primitive.Add,
		Value: int(value),
	}
}

// Decrement :
func Decrement[T constraints.Integer](field string, value T) primitive.Math {
	return primitive.Math{
		Field: field,
		Mode:  primitive.Deduct,
		Value: int(value),
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
