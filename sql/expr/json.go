package expr

import (
	"encoding/json"
	"github.com/si3nloong/sqlike/sqlike/primitive"
)

// JSON_QUOTE :
func JSON_QUOTE(doc interface{}) (f primitive.JSONFunc) {
	f.Type = primitive.JSONQuote
	switch vi := doc.(type) {
	case string:
		f.Args = append(f.Args, primitive.Value{
			Raw: vi,
		})
	default:
		f.Args = append(f.Args, vi)
	}
	return
}

// JSON_UNQUOTE :
func JSON_UNQUOTE(doc interface{}) (f primitive.JSONFunc) {
	f.Type = primitive.JSONUnquote
	switch vi := doc.(type) {
	case string:
		f.Args = append(f.Args, primitive.Value{
			Raw: vi,
		})
	default:
		f.Args = append(f.Args, vi)
	}
	return
}

// JSON_EXTRACT :
func JSON_EXTRACT(doc interface{}, path string, otherPaths ...string) (fc primitive.JSONFunc) {
	fc.Type = primitive.JSONExtract
	fc.Args = append(fc.Args, doc)
	for _, p := range append([]string{path}, otherPaths...) {
		fc.Args = append(fc.Args, primitive.Value{
			Raw: p,
		})
	}
	return
}

// JSON_KEYS :
func JSON_KEYS(doc interface{}, paths ...string) (fc primitive.JSONFunc) {
	fc.Type = primitive.JSONKeys
	fc.Args = append(fc.Args, doc)
	for _, p := range paths {
		fc.Args = append(fc.Args, primitive.Value{
			Raw: p,
		})
	}
	return
}

// JSON_VALID :
func JSON_VALID(val interface{}) (fc primitive.JSONFunc) {
	fc.Type = primitive.JSONValid
	fc.Args = append(fc.Args, val)
	return
}

// JSON_CONTAINS :
func JSON_CONTAINS(target, candidate interface{}, paths ...string) (f primitive.JSONFunc) {
	f.Type = primitive.JSONContains
	for _, arg := range []interface{}{target, candidate} {
		switch vi := arg.(type) {
		case string, json.RawMessage:
			f.Args = append(f.Args, primitive.Value{
				Raw: vi,
			})
		case primitive.Column:
			f.Args = append(f.Args, vi)
		default:
			f.Args = append(f.Args, vi)
		}
	}
	if len(paths) > 0 {
		for _, p := range paths {
			f.Args = append(f.Args, p)
		}
	}
	return
}

// func JSON_REPLACE()
// func JSON_TYPE()
// func JSON_VALID()

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
