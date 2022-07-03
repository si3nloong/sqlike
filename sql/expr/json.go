package expr

import (
	"encoding/json"

	"github.com/si3nloong/sqlike/v2/x/primitive"
)

// JSON_QUOTE :
// SELECT JSON_QUOTE(`Column` -> '$.type') FROM test;
func JSON_QUOTE(doc any) (f primitive.JSONFunc) {
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
func JSON_UNQUOTE(doc any) (f primitive.JSONFunc) {
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
func JSON_EXTRACT(doc any, path string, otherPaths ...string) (f primitive.JSONFunc) {
	f.Type = primitive.JSON_EXTRACT
	f.Args = append(f.Args, doc)
	for _, p := range append([]string{path}, otherPaths...) {
		f.Args = append(f.Args, wrapRaw(p))
	}
	return
}

// JSON_KEYS :
func JSON_KEYS(doc any, paths ...string) (f primitive.JSONFunc) {
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
func JSON_SET(doc any, path string, value any, pathValues ...any) (f primitive.JSONFunc) {
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
func JSON_INSERT(doc any, path string, value any, pathValues ...any) (f primitive.JSONFunc) {
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
func JSON_REMOVE(doc any, path string, paths ...string) (f primitive.JSONFunc) {
	f.Type = primitive.JSON_REMOVE
	f.Args = append(f.Args, doc, wrapRaw(path))
	for _, p := range paths {
		f.Args = append(f.Args, wrapRaw(p))
	}
	return
}

// JSON_REPLACE :
func JSON_REPLACE(doc any, path string, value any, pathValues ...any) (f primitive.JSONFunc) {
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
func JSON_VALID(val any) (f primitive.JSONFunc) {
	f.Type = primitive.JSON_VALID
	f.Args = append(f.Args, val)
	return
}

// JSON_CONTAINS :
func JSON_CONTAINS(target, candidate any, paths ...string) (f primitive.JSONFunc) {
	f.Type = primitive.JSON_CONTAINS
	for _, arg := range []any{target, candidate} {
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
func JSON_TYPE(val any) (f primitive.JSONFunc) {
	f.Type = primitive.JSON_TYPE
	f.Args = append(f.Args, val)
	return
}

// MemberOf : mysql 8.0.17
func MemberOf(val any, arr any) (f primitive.JSONFunc) {
	f.Prefix = val
	f.Type = primitive.MEMBER_OF
	f.Args = append(f.Args, wrapColumn(arr))
	return
}

// JSONColumn :
func JSONColumn(column string, nested ...string) (c primitive.JSONColumn) {
	c.Column = column
	c.Nested = nested
	c.UnquoteResult = false
	return
}
