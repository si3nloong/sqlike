package expr

import (
	"github.com/si3nloong/sqlike/sqlike/primitive"
)

// Sum :
func Sum(field interface{}) (a primitive.Aggregate) {
	a.Field = wrapColumn(field)
	a.By = primitive.Sum
	return
}

// Count :
func Count(field interface{}) (a primitive.Aggregate) {
	a.Field = wrapColumn(field)
	a.By = primitive.Count
	return
}

// Average :
func Average(field interface{}) (a primitive.Aggregate) {
	a.Field = wrapColumn(field)
	a.By = primitive.Average
	return
}

// Max :
func Max(field interface{}) (a primitive.Aggregate) {
	a.Field = wrapColumn(field)
	a.By = primitive.Max
	return
}

// Min :
func Min(field interface{}) (a primitive.Aggregate) {
	a.Field = wrapColumn(field)
	a.By = primitive.Min
	return
}
