package expr

import (
	"reflect"

	"github.com/si3nloong/sqlike/sqlike/primitive"
)

// Equal :
func Equal(field, value interface{}) (c primitive.C) {
	c = clause(field, primitive.Equal, value)
	return
}

// NotEqual :
func NotEqual(field, value interface{}) (c primitive.C) {
	c = clause(field, primitive.NotEqual, value)
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
func In(field, values interface{}) (c primitive.C) {
	c = inGroup(field, primitive.In, values)
	return
}

// NotIn :
func NotIn(field, values interface{}) (c primitive.C) {
	c = inGroup(field, primitive.NotIn, values)
	return
}

func inGroup(field interface{}, op primitive.Operator, values interface{}) (c primitive.C) {
	v := reflect.ValueOf(values)
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
		grp = append(grp, values)
	}
	grp = append(grp, Raw(")"))
	c.Value = grp
	return c
}

// Like :
func Like(field, value interface{}) (c primitive.C) {
	c = clause(field, primitive.Like, value)
	return
}

// NotLike :
func NotLike(field, value interface{}) (c primitive.C) {
	c = clause(field, primitive.NotLike, value)
	return
}

// GreaterOrEqual :
func GreaterOrEqual(field, value interface{}) (c primitive.C) {
	c = clause(field, primitive.GreaterOrEqual, value)
	return
}

// GreaterThan :
func GreaterThan(field, value interface{}) (c primitive.C) {
	c = clause(field, primitive.GreaterThan, value)
	return
}

// LesserOrEqual :
func LesserOrEqual(field, value interface{}) (c primitive.C) {
	c = clause(field, primitive.LesserOrEqual, value)
	return
}

// LesserThan :
func LesserThan(field, value interface{}) (c primitive.C) {
	c = clause(field, primitive.LesserThan, value)
	return
}

// Between :
func Between(field, from, to interface{}) (c primitive.C) {
	c = clause(field, primitive.Between, primitive.R{From: from, To: to})
	return
}

// NotBetween :
func NotBetween(field, from, to interface{}) (c primitive.C) {
	c = clause(field, primitive.NotBetween, primitive.R{From: from, To: to})
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

func clause(field interface{}, op primitive.Operator, value interface{}) (c primitive.C) {
	c.Field = wrapColumn(field)
	c.Operator = op
	c.Value = value
	return
}
