package expr

import (
	"github.com/si3nloong/sqlike/sqlike/primitive"
)

// JSONQuote :
func JSON_QUOTE(value string) (fc primitive.JSONFunc) {
	fc.Type = primitive.JSONQuote
	fc.Arguments = append(fc.Arguments, wrapColumn(value))
	return
}

// JSON_CONTAINS :
func JSON_CONTAINS(target, candidate interface{}, paths ...string) (jc primitive.JC) {
	var path *string
	if len(paths) > 0 {
		path = &paths[0]
	}
	switch vi := target.(type) {
	case string:
		jc.Target = primitive.Column{Name: vi}
	case primitive.Column:
		jc.Target = vi
	default:
		jc.Target = primitive.Value{Raw: vi}
	}
	jc.Candidate = wrapJSONColumn(candidate)
	jc.Path = path
	return
}

// func JSON_REPLACE()
// func JSON_TYPE()
// func JSON_VALID()
// func JSON_UNQUOTE()

// JSONColumn :
func JSONColumn(column string, nested ...string) (c primitive.JSONColumn) {
	c.Column = column
	c.Nested = nested
	c.UnquoteResult = false
	return
}

func wrapJSONColumn(it interface{}) interface{} {
	switch vi := it.(type) {
	case primitive.Column:
		return primitive.CastAs{
			Value:    vi,
			DataType: primitive.JSON,
		}
	case primitive.JSONFunc:
		return vi
	default:
		return primitive.Value{Raw: vi}
	}
}
