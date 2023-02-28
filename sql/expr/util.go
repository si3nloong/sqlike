package expr

import "github.com/si3nloong/sqlike/v2/internal/primitive"

func wrapColumns(v any) primitive.Pair {
	switch vi := v.(type) {
	case string:
		return Pair("", vi)
	case primitive.Pair:
		return Pair(vi[0], vi[1])
	default:
		panic("sqlike: invalid type to wrap")
	}
}
