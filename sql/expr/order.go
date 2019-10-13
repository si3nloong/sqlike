package expr

import (
	"fmt"
	"reflect"

	"github.com/si3nloong/sqlike/sqlike/primitive"
)

// Field :
func Field(name string, val interface{}) (f primitive.Field) {
	f.Name = name
	v := reflect.ValueOf(val)
	k := v.Kind()
	if k != reflect.Array && k != reflect.Slice {
		panic(fmt.Errorf("unsupported data type: %v", k))
	}
	length := v.Len()
	if length < 1 {
		panic("zero length of array/slice")
	}
	for i := 0; i < length; i++ {
		f.Values = append(f.Values, v.Index(i).Interface())
	}
	return
}

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
