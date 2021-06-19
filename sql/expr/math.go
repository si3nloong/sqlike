package expr

import "github.com/si3nloong/sqlike/x/primitive"

// Increment :
func Increment(field string, inc uint) primitive.Math {
	return primitive.Math{
		Field: field,
		Mode:  primitive.Add,
		Value: int(inc),
	}
}

// Decrement :
func Decrement(field string, inc uint) primitive.Math {
	return primitive.Math{
		Field: field,
		Mode:  primitive.Deduct,
		Value: int(inc),
	}
}

// Multiply :
func Multiply(fields ...interface{}) (grp primitive.Group) {
	for i, f := range fields {
		if i > 0 {
			grp.Values = append(grp.Values, Raw("*"))
		}
		grp.Values = append(grp.Values, wrapColumn(f))
	}
	return
}
