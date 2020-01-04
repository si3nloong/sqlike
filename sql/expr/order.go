package expr

import (
	"fmt"
	"reflect"

	"github.com/si3nloong/sqlike/reflext"
	"github.com/si3nloong/sqlike/sqlike/primitive"
)

// Field :
func Field(name string, val interface{}) (f primitive.Field) {
	f.Name = name
	v := reflext.Indirect(reflect.ValueOf(val))
	k := v.Kind()
	if k != reflect.Array && k != reflect.Slice {
		panic(fmt.Errorf("unsupported data type: %v", k))
	}
	length := v.Len()
	if length < 1 {
		panic("zero length of array or slice")
	}
	for i := 0; i < length; i++ {
		f.Values = append(f.Values, v.Index(i).Interface())
	}
	return
}

// Asc :
func Asc(field interface{}) (s primitive.Sort) {
	s.Field = wrapColumn(field)
	s.Order = primitive.Ascending
	return
}

// Desc :
func Desc(field interface{}) (s primitive.Sort) {
	s.Field = wrapColumn(field)
	s.Order = primitive.Descending
	return
}
