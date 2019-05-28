package expr

import "bitbucket.org/SianLoong/sqlike/sqlike/primitive"

// Increment :
func Increment(field string, inc uint) primitive.Math {
	return primitive.Math{
		Field: primitive.Col(field),
		Mode:  primitive.Add,
		Value: int(inc),
	}
}

// Inc :
func Inc(field string, inc uint) primitive.Math {
	return primitive.Math{
		Field: primitive.Col(field),
		Mode:  primitive.Add,
		Value: int(inc),
	}
}
