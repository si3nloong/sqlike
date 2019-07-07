package expr

import (
	"fmt"
	"reflect"

	"github.com/si3nloong/sqlike/sqlike/primitive"
)

// Raw :
func Raw(value string) (r primitive.Raw) {
	r.Value = value
	return
}

// Column :
func Column(name string) (c primitive.Column) {
	c.Name = name
	return
}

// Equal :
func Equal(field, value interface{}) (c primitive.C) {
	c.Field = wrapColumn(field)
	c.Operator = primitive.Equal
	c.Value = value
	return
}

// NotEqual :
func NotEqual(field, value interface{}) (c primitive.C) {
	c.Field = wrapColumn(field)
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
func In(field, value interface{}) (c primitive.C) {
	v := reflect.ValueOf(value)
	k := v.Kind()
	c.Field = wrapColumn(field)
	c.Operator = primitive.In
	grp := primitive.G{}
	grp = append(grp, Raw("("))
	if k == reflect.Array || k == reflect.Slice {
		for i := 0; i < v.Len(); i++ {
			if i > 0 {
				grp = append(grp, Raw(","))
			}
			grp = append(grp, v.Index(i).Interface())
		}
	} else {
		grp = append(grp, value)
	}
	grp = append(grp, Raw(")"))
	c.Value = grp
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
	grp := primitive.G{}
	for i := 0; i < v.Len(); i++ {
		grp = append(grp, v.Index(i).Interface())
	}
	c.Value = grp
	return
}

// Like :
func Like(field, value interface{}) (c primitive.C) {
	c.Field = wrapColumn(field)
	c.Operator = primitive.Like
	c.Value = value
	return
}

// NotLike :
func NotLike(field, value interface{}) (c primitive.C) {
	c.Field = wrapColumn(field)
	c.Operator = primitive.NotLike
	c.Value = value
	return
}

// GreaterOrEqual :
func GreaterOrEqual(field, value interface{}) (c primitive.C) {
	c.Field = wrapColumn(field)
	c.Operator = primitive.GreaterEqual
	c.Value = value
	return
}

// GreaterThan :
func GreaterThan(field, value interface{}) (c primitive.C) {
	c.Field = wrapColumn(field)
	c.Operator = primitive.GreaterThan
	c.Value = value
	return
}

// LesserOrEqual :
func LesserOrEqual(field, value interface{}) (c primitive.C) {
	c.Field = wrapColumn(field)
	c.Operator = primitive.LowerEqual
	c.Value = value
	return
}

// LesserThan :
func LesserThan(field, value interface{}) (c primitive.C) {
	c.Field = wrapColumn(field)
	c.Operator = primitive.LowerThan
	c.Value = value
	return
}

// Between :
func Between(field, from, to interface{}) (c primitive.C) {
	c.Field = wrapColumn(field)
	c.Operator = primitive.Between
	c.Value = primitive.R{From: from, To: to}
	return
}

// NotBetween :
func NotBetween(field, from, to interface{}) (c primitive.C) {
	c.Field = wrapColumn(field)
	c.Operator = primitive.NotBetween
	c.Value = primitive.R{From: from, To: to}
	return
}

// And :
func And(conds ...interface{}) (g primitive.G) {
	if len(conds) > 1 {
		g = append(g, Raw("("))
		for i, cond := range conds {
			if i > 0 {
				g = append(g, primitive.And)
			}
			g = append(g, cond)
		}
		g = append(g, Raw(")"))
		return
	}
	g = append(g, conds...)
	return
}

// Or :
func Or(conds ...interface{}) (g primitive.G) {
	if len(conds) > 1 {
		g = append(g, Raw("("))
		for i, cond := range conds {
			if i > 0 {
				g = append(g, primitive.Or)
			}
			g = append(g, Raw("("), cond, Raw(")"))
		}
		g = append(g, Raw(")"))
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

// Date :
func Date(field string) (d primitive.Func) {
	return
}

// Asc :
func Asc(field string) (s primitive.Sort) {
	s.Field = field
	s.Order = primitive.Ascending
	return
}

func wrapColumn(it interface{}) interface{} {
	switch vi := it.(type) {
	case string:
		return Column(vi)
	case primitive.Column:
		return vi
	default:
		return vi
	}
}

// Desc :
func Desc(field string) (s primitive.Sort) {
	s.Field = field
	s.Order = primitive.Descending
	return
}
