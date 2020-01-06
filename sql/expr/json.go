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
func JSON_EXTRACT(doc interface{}, path string, otherPaths ...string) (f primitive.JSONFunc) {
	f.Type = primitive.JSONExtract
	f.Args = append(f.Args, doc)
	for _, p := range append([]string{path}, otherPaths...) {
		f.Args = append(f.Args, primitive.Value{
			Raw: p,
		})
	}
	return
}

// JSON_KEYS :
func JSON_KEYS(doc interface{}, paths ...string) (f primitive.JSONFunc) {
	f.Type = primitive.JSONKeys
	f.Args = append(f.Args, doc)
	for _, p := range paths {
		f.Args = append(f.Args, primitive.Value{
			Raw: p,
		})
	}
	return
}

// JSON_VALID :
func JSON_VALID(val interface{}) (f primitive.JSONFunc) {
	f.Type = primitive.JSONValid
	f.Args = append(f.Args, val)
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

// JSON_TYPE :
func JSON_TYPE(val interface{}) (f primitive.JSONFunc) {
	// f.Type = primitive.JSONType
	f.Args = append(f.Args, val)
	return
}

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
