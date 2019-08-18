package expr

import "github.com/si3nloong/sqlike/sqlike/primitive"

// Increment :
func Increment(field string, inc uint) primitive.Math {
	return primitive.Math{
		Field: primitive.Col(field),
		Mode:  primitive.Add,
		Value: int(inc),
	}
}

// Multiply :
func Multiply(fields ...interface{}) (grp primitive.Group) {
	for i, f := range fields {
		if i > 0 {
			grp = append(grp, Raw("*"))
		}
		grp = append(grp, wrapColumn(f))
	}
	return
}
