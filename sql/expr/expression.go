package expr

import (
	"reflect"

	"github.com/si3nloong/sqlike/reflext"
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
func IsNull(field string) (c primitive.Nil) {
	c.Field = wrapColumn(field)
	c.IsNot = true
	return
}

// NotNull :
func NotNull(field string) (c primitive.Nil) {
	c.Field = wrapColumn(field)
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

// Like :
func Like(field, value interface{}) (p primitive.L) {
	p.Field = wrapColumn(field)
	p.Value = value
	return
}

// NotLike :
func NotLike(field, value interface{}) (p primitive.L) {
	p.Field = wrapColumn(field)
	p.IsNot = true
	p.Value = value
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
func And(conds ...interface{}) (g primitive.Group) {
	g = buildGroup(primitive.And, conds)
	return
}

// Or :
func Or(conds ...interface{}) (g primitive.Group) {
	g = buildGroup(primitive.Or, conds)
	return
}

func buildGroup(op primitive.Operator, conds []interface{}) (g primitive.Group) {
	length := len(conds)
	if length < 1 {
		return
	}
	if length == 1 {
		x, ok := conds[0].(primitive.Group)
		if ok {
			g = x
			return
		}
	}

	sg := make([]interface{}, 0, length)
	for len(conds) > 0 {
		cond := conds[0]
		conds = conds[1:]

		if cond == nil || reflext.IsZero(reflext.ValueOf(cond)) {
			continue
		}

		if len(sg) > 0 {
			sg = append(sg, op)
		}
		sg = append(sg, cond)
	}
	if len(sg) > 1 {
		g.Values = append(g.Values, Raw("("))
		g.Values = append(g.Values, sg...)
		g.Values = append(g.Values, Raw(")"))
		return
	}
	g.Values = append(g.Values, sg...)
	return
}

// ColumnValue :
func ColumnValue(field string, value interface{}) (kv primitive.KV) {
	kv.Field = field
	kv.Value = value
	return
}

// CastAs :
func CastAs(value interface{}, datatype primitive.DataType) (cast primitive.CastAs) {
	cast.Value = value
	cast.DataType = datatype
	return
}

func inGroup(field interface{}, op primitive.Operator, values interface{}) (c primitive.C) {
	v := reflect.ValueOf(values)
	k := v.Kind()
	c.Field = wrapColumn(field)
	c.Operator = primitive.In
	grp := primitive.Group{}
	grp.Values = append(grp.Values, Raw("("))
	if k == reflect.Array || k == reflect.Slice {
		for i := 0; i < v.Len(); i++ {
			if i > 0 {
				grp.Values = append(grp.Values, Raw(","))
			}
			grp.Values = append(grp.Values, v.Index(i).Interface())
		}
	} else {
		grp.Values = append(grp.Values, values)
	}
	grp.Values = append(grp.Values, Raw(")"))
	c.Value = grp
	return c
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
