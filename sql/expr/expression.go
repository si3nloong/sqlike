package expr

import (
	"reflect"

	"github.com/si3nloong/sqlike/v2/internal/primitive"
	"github.com/si3nloong/sqlike/v2/x/reflext"
)

type ColumnConstraints interface {
	~string | primitive.Column | primitive.Pair | primitive.JSONColumn
}

// Equal :
func Equal[C ColumnConstraints](field C, value any) (c primitive.C) {
	c = clause(field, primitive.Equal, value)
	return
}

// NotEqual :
func NotEqual[C ColumnConstraints](field C, value any) (c primitive.C) {
	c = clause(field, primitive.NotEqual, value)
	return
}

// IsNull :
func IsNull[C ColumnConstraints](field C) (c primitive.Nil) {
	c.Field = wrapColumn(field)
	return
}

// IsNotNull :
func IsNotNull[C ColumnConstraints](field C) (c primitive.Nil) {
	c.Field = wrapColumn(field)
	c.IsNot = true
	return
}

// In :
func In[C ColumnConstraints](field C, values any) (c primitive.C) {
	c = inGroup(field, primitive.In, values)
	return
}

// NotIn :
func NotIn[C ColumnConstraints](field C, values any) (c primitive.C) {
	c = inGroup(field, primitive.NotIn, values)
	return
}

// Any :
func Any[T any](values []T) (c primitive.C) {
	return
}

// Like :
func Like[F ColumnConstraints, V string | primitive.Raw](field F, value V) (p primitive.L) {
	p.Field = wrapColumn(field)
	p.Value = value
	return
}

// NotLike :
func NotLike[C ColumnConstraints](field C, value any) (p primitive.L) {
	p.Field = wrapColumn(field)
	p.IsNot = true
	p.Value = value
	return
}

// GreaterThan :
func GreaterThan(field, value any) (c primitive.C) {
	c = clause(field, primitive.GreaterThan, value)
	return
}

// GreaterOrEqual :
func GreaterOrEqual(field, value any) (c primitive.C) {
	c = clause(field, primitive.GreaterOrEqual, value)
	return
}

// LesserThan :
func LesserThan(field, value any) (c primitive.C) {
	c = clause(field, primitive.LesserThan, value)
	return
}

// LesserOrEqual :
func LesserOrEqual(field, value any) (c primitive.C) {
	c = clause(field, primitive.LesserOrEqual, value)
	return
}

// Between :
func Between(field, from, to any) (c primitive.C) {
	c = clause(field, primitive.Between, primitive.R{From: from, To: to})
	return
}

// NotBetween :
func NotBetween(field, from, to any) (c primitive.C) {
	c = clause(field, primitive.NotBetween, primitive.R{From: from, To: to})
	return
}

// And :
func And(conds ...any) (g primitive.Group) {
	g = buildGroup(primitive.And, conds)
	return
}

// Or :
func Or(conds ...any) (g primitive.Group) {
	g = buildGroup(primitive.Or, conds)
	return
}

func buildGroup(op primitive.Operator, conds []any) (g primitive.Group) {
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

	sg := make([]any, 0, length)
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
func ColumnValue[C ColumnConstraints](field C, value any) (kv primitive.KV) {
	switch vi := any(field).(type) {
	case string:
		kv.Field = vi
	case primitive.Column:
		kv.Field = vi.Name
	case primitive.Pair:
		kv.Field = vi[1]
	}
	kv.Value = value
	return
}

// CastAs :
func CastAs(value any, datatype primitive.DataType) (cast primitive.CastAs) {
	cast.Value = value
	cast.DataType = datatype
	return
}

func inGroup(field any, op primitive.Operator, values any) (c primitive.C) {
	v := reflect.ValueOf(values)
	k := v.Kind()
	c.Field = wrapColumn(field)
	c.Operator = op
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

func wrapColumn(it any) any {
	switch vi := it.(type) {
	case string:
		return Column(vi)
	case primitive.Pair:
		return Column(vi[0], vi[1])
	case primitive.Column:
		return vi
	default:
		return vi
	}
}

func clause(field any, op primitive.Operator, value any) (c primitive.C) {
	c.Field = wrapColumn(field)
	c.Operator = op
	c.Value = value
	return
}
