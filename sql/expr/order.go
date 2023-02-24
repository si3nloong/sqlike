package expr

import (
	"github.com/si3nloong/sqlike/v2/internal/primitive"
)

// Field :
func Field[T any](name string, val []T) (f primitive.Field) {
	f.Name = name
	length := len(val)
	if length < 1 {
		panic("zero length of array or slice")
	}
	for i := 0; i < length; i++ {
		f.Values = append(f.Values, val[i])
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
