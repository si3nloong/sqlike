package expr

import (
	"encoding/json"

	"github.com/si3nloong/sqlike/sqlike/primitive"
)

// JSON_QUOTE :
func JSON_QUOTE(doc interface{}) (f primitive.JSONFunc) {
	f.Type = primitive.JSON_QUOTE
	switch vi := doc.(type) {
	case string:
		f.Args = append(f.Args, wrapRaw(vi))
	default:
		f.Args = append(f.Args, vi)
	}
	return
}

// JSON_UNQUOTE :
func JSON_UNQUOTE(doc interface{}) (f primitive.JSONFunc) {
	f.Type = primitive.JSON_UNQUOTE
	switch vi := doc.(type) {
	case string:
		f.Args = append(f.Args, wrapRaw(vi))
	default:
		f.Args = append(f.Args, vi)
	}
	return
}

// JSON_EXTRACT :
func JSON_EXTRACT(doc interface{}, path string, otherPaths ...string) (f primitive.JSONFunc) {
	f.Type = primitive.JSON_EXTRACT
	f.Args = append(f.Args, doc)
	for _, p := range append([]string{path}, otherPaths...) {
		f.Args = append(f.Args, wrapRaw(p))
	}
	return
}

// JSON_KEYS :
func JSON_KEYS(doc interface{}, paths ...string) (f primitive.JSONFunc) {
	f.Type = primitive.JSON_KEYS
	f.Args = append(f.Args, doc)
	for _, p := range paths {
		f.Args = append(f.Args, primitive.Value{
			Raw: p,
		})
	}
	return
}

// JSON_SET :
func JSON_SET(doc interface{}, path string, value interface{}, pathValues ...interface{}) (f primitive.JSONFunc) {
	length := len(pathValues)
	if length > 0 && length%2 != 0 {
		panic("invalid argument len for JSON_SET(json_doc, path, val[, path, val] ...)")
	}
	f.Type = primitive.JSON_SET
	f.Args = append(f.Args, doc, wrapRaw(path), value)
	f.Args = append(f.Args, pathValues...)
	return
}

// JSON_INSERT :
func JSON_INSERT(doc interface{}, path string, value interface{}, pathValues ...interface{}) (f primitive.JSONFunc) {
	length := len(pathValues)
	if length > 0 && length%2 != 0 {
		panic("invalid argument len for JSON_INSERT(json_doc, path, val[, path, val] ...)")
	}
	f.Type = primitive.JSON_INSERT
	f.Args = append(f.Args, doc, wrapRaw(path), value)
	f.Args = append(f.Args, pathValues...)
	return
}

// JSON_REMOVE :
func JSON_REMOVE(doc interface{}, path string, paths ...string) (f primitive.JSONFunc) {
	f.Type = primitive.JSON_REMOVE
	f.Args = append(f.Args, doc, wrapRaw(path))
	for _, p := range paths {
		f.Args = append(f.Args, wrapRaw(p))
	}
	return
}

// JSON_REPLACE :
func JSON_REPLACE(doc interface{}, path string, value interface{}, pathValues ...interface{}) (f primitive.JSONFunc) {
	length := len(pathValues)
	if length > 0 && length%2 != 0 {
		panic("invalid argument len for JSON_REPLACE(json_doc, path, val[, path, val] ...)")
	}
	f.Type = primitive.JSON_REPLACE
	f.Args = append(f.Args, doc, wrapRaw(path), value)
	f.Args = append(f.Args, pathValues...)
	return
}

// JSON_VALID :
func JSON_VALID(val interface{}) (f primitive.JSONFunc) {
	f.Type = primitive.JSON_VALID
	f.Args = append(f.Args, val)
	return
}

// JSON_CONTAINS :
func JSON_CONTAINS(target, candidate interface{}, paths ...string) (f primitive.JSONFunc) {
	f.Type = primitive.JSON_CONTAINS
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
	f.Type = primitive.JSON_TYPE
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
