package expr

import (
	"github.com/si3nloong/sqlike/sqlike/primitive"
)

// JSON_QUOTE :
func JSON_QUOTE(doc string) (fc primitive.JSONFunc) {
	fc.Type = primitive.JSONQuote
	fc.Arguments = append(fc.Arguments, wrapColumn(doc))
	return
}

// JSON_UNQUOTE :
func JSON_UNQUOTE(doc interface{}) (fc primitive.JSONFunc) {
	fc.Type = primitive.JSONUnquote
	fc.Arguments = append(fc.Arguments, wrapColumn(doc))
	return
}

// JSON_EXTRACT :
func JSON_EXTRACT(doc interface{}, path string, otherPaths ...string) (fc primitive.JSONFunc) {
	fc.Type = primitive.JSONExtract
	fc.Arguments = append(fc.Arguments, doc)
	for _, p := range append([]string{path}, otherPaths...) {
		fc.Arguments = append(fc.Arguments, primitive.Value{
			Raw: p,
		})
	}
	return
}

// JSON_KEYS :
func JSON_KEYS(doc interface{}, paths ...string) (fc primitive.JSONFunc) {
	fc.Type = primitive.JSONKeys
	fc.Arguments = append(fc.Arguments, doc)
	for _, p := range paths {
		fc.Arguments = append(fc.Arguments, primitive.Value{
			Raw: p,
		})
	}
	return
}

// JSON_VALID :
func JSON_VALID(val interface{}) (fc primitive.JSONFunc) {
	fc.Type = primitive.JSONValid
	fc.Arguments = append(fc.Arguments, val)
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
