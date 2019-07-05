package expr

import (
	"fmt"
	"reflect"

	"github.com/si3nloong/sqlike/sqlike/primitive"
)

// Equal :
func Equal(field string, value interface{}) (c primitive.C) {
	c.Field = primitive.L(field)
	c.Operator = primitive.Equal
	c.Value = value
	return
}

// NotEqual :
func NotEqual(field string, value interface{}) (c primitive.C) {
	c.Field = primitive.L(field)
	c.Operator = primitive.NotEqual
	c.Value = value
	return
}

// IsNull :
func IsNull(field string) (c primitive.C) {
	c.Field = primitive.L(field)
	c.Operator = primitive.IsNull
	return
}

// NotNull :
func NotNull(field string) (c primitive.C) {
	c.Field = primitive.L(field)
	c.Operator = primitive.NotNull
	return
}

// In :
func In(field string, value interface{}) (c primitive.C) {
	v := reflect.ValueOf(value)
	k := v.Kind()
	c.Field = primitive.L(field)
	c.Operator = primitive.In
	if k == reflect.Array || k == reflect.Slice {
		grp := primitive.GV{}
		for i := 0; i < v.Len(); i++ {
			grp = append(grp, v.Index(i).Interface())
		}
		c.Value = grp
		return c
	}
	c.Value = value
	return c
}

// NotIn :
func NotIn(field string, value interface{}) (c primitive.C) {
	v := reflect.ValueOf(value)
	c.Field = primitive.L(field)
	c.Operator = primitive.NotIn
	k := v.Kind()
	if k != reflect.Array && k != reflect.Slice {
		panic(fmt.Errorf("expr.In not support data type %v", v.Type()))
	}
	grp := primitive.GV{}
	for i := 0; i < v.Len(); i++ {
		grp = append(grp, v.Index(i).Interface())
	}
	c.Value = grp
	return
}

// Like :
func Like(field string, value interface{}) (c primitive.C) {
	c.Field = primitive.L(field)
	c.Operator = primitive.Like
	c.Value = value
	return
}

// NotLike :
func NotLike(field string, value interface{}) (c primitive.C) {
	c.Field = primitive.L(field)
	c.Operator = primitive.NotLike
	c.Value = value
	return
}

// GreaterEqual :
func GreaterEqual(field string, value interface{}) (c primitive.C) {
	c.Field = primitive.L(field)
	c.Operator = primitive.GreaterEqual
	c.Value = value
	return
}

// GreaterThan :
func GreaterThan(field string, value interface{}) (c primitive.C) {
	c.Field = primitive.L(field)
	c.Operator = primitive.GreaterThan
	c.Value = value
	return
}

// LowerEqual :
func LowerEqual(field string, value interface{}) (c primitive.C) {
	c.Field = primitive.L(field)
	c.Operator = primitive.LowerEqual
	c.Value = value
	return
}

// LowerThan :
func LowerThan(field string, value interface{}) (c primitive.C) {
	c.Field = primitive.L(field)
	c.Operator = primitive.LowerThan
	c.Value = value
	return
}

// Between :
func Between(field string, from, to interface{}) (c primitive.C) {
	c.Field = primitive.L(field)
	c.Operator = primitive.Between
	c.Value = primitive.R{From: from, To: to}
	return
}

// NotBetween :
func NotBetween(field string, from, to interface{}) (c primitive.C) {
	c.Field = primitive.L(field)
	c.Operator = primitive.NotBetween
	c.Value = primitive.R{From: from, To: to}
	return
}

// And :
func And(conds ...interface{}) (g primitive.G) {
	if len(conds) > 1 {
		g = append(g, primitive.Raw(`(`))
		for i, cond := range conds {
			if i > 0 {
				g = append(g, primitive.And)
			}
			g = append(g, cond)
		}
		g = append(g, primitive.Raw(`)`))
		return
	}
	g = append(g, conds...)
	return
}

// Or :
func Or(conds ...interface{}) (g primitive.G) {
	if len(conds) > 1 {
		g = append(g, primitive.Raw(`(`))
		for i, cond := range conds {
			if i > 0 {
				g = append(g, primitive.Or)
			}
			g = append(g, primitive.Raw(`(`), cond, primitive.Raw(`)`))
		}
		g = append(g, primitive.Raw(`)`))
		return
	}
	g = append(g, conds...)
	return
}

// Field :
func Field(field string, value interface{}) (kv primitive.KV) {
	kv.Field = field
	kv.Value = value
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
