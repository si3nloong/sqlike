package expr

import "github.com/si3nloong/sqlike/sqlike/primitive"

// Asc :
func Asc(field string) (s primitive.Sort) {
	s.Field = field
	s.Order = primitive.Ascending
	return
}

// Desc :
func Desc(field string) (s primitive.Sort) {
	s.Field = field
	s.Order = primitive.Descending
	return
}
