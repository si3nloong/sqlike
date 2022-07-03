package expr

import (
	"github.com/si3nloong/sqlike/v2/x/primitive"
)

// Sum :
func Sum[T string | int](field T) (a primitive.Aggregate) {
	a.Field = wrapColumn(field)
	a.By = primitive.Sum
	return
}

// Count :
func Count(field any) (a primitive.Aggregate) {
	a.Field = wrapColumn(field)
	a.By = primitive.Count
	return
}

// Average :
func Average(field any) (a primitive.Aggregate) {
	a.Field = wrapColumn(field)
	a.By = primitive.Average
	return
}

// Max :
func Max(field any) (a primitive.Aggregate) {
	a.Field = wrapColumn(field)
	a.By = primitive.Max
	return
}

// Min :
func Min(field any) (a primitive.Aggregate) {
	a.Field = wrapColumn(field)
	a.By = primitive.Min
	return
}
