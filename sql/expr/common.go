package expr

import "github.com/si3nloong/sqlike/sqlike/primitive"

// Raw :
func Raw(value string) (r primitive.Raw) {
	r.Value = value
	return
}

// As :
func As(field interface{}, alias string) (as primitive.As) {
	as.Field = wrapColumn(field)
	as.Name = alias
	return
}

// Column :
func Column(name string, alias ...string) (c primitive.Column) {
	if len(alias) > 0 {
		c.Table = alias[0]
	}
	c.Name = name
	return
}

// JSONColumn :
func JSONColumn(column string, nested ...string) (c primitive.JSONColumn) {
	c.Column = column
	c.Nested = nested
	c.UnquoteResult = false
	return
}
