package expr

import (
	"github.com/si3nloong/sqlike/v2/internal/primitive"
)

// Field :
func Field[T any](name string, val []T) (f primitive.Field) {
	f.Name = name
	if len(val) < 1 {
		panic(`sqlike: zero length of array or slice`)
	}
	for _, v := range val {
		f.Values = append(f.Values, v)
	}
	return
}

// Asc :
func Asc[C ColumnConstraints](field C) (s primitive.Sort) {
	s.Field = wrapColumn(field)
	s.Order = primitive.Ascending
	return
}

// Desc :
func Desc[C ColumnConstraints](field C) (s primitive.Sort) {
	s.Field = wrapColumn(field)
	s.Order = primitive.Descending
	return
}
