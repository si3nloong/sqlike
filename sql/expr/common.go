package expr

import (
	"database/sql"
	"strings"

	"github.com/si3nloong/sqlike/v2/internal/primitive"
)

func Pair(first, second string) (p primitive.Pair) {
	p = [2]string{first, second}
	return
}

// Raw :
func Raw(value string) (r primitive.Raw) {
	r.Value = value
	return
}

// As :
func As(src any, alias string) (as primitive.As) {
	as.Field = wrapColumn(src)
	as.Name = alias
	return
}

// Column :
func Column(name string, alias ...string) (c primitive.Column) {
	if len(alias) > 0 {
		c.Table = name
		c.Name = alias[0]
		return
	}
	c.Name = name
	return
}

// Func :
func Func(name string, value any, others ...any) (f primitive.Func) {
	f.Name = strings.ToUpper(strings.TrimSpace(name))
	f.Args = append(f.Args, wrapRaw(value))
	if len(others) > 0 {
		for _, arg := range others {
			f.Args = append(f.Args, wrapRaw(arg))
		}
	}
	return
}

func wrapRaw(v any) (it any) {
	vv := primitive.Value{}
	switch vi := v.(type) {
	case sql.RawBytes:
		vv.Raw = vi
		return vv
	case nil:
		vv.Raw = vi
		return vv
	case string:
		vv.Raw = vi
		return vv
	case []byte:
		vv.Raw = vi
		return vv
	case float32, float64:
		vv.Raw = vi
		return vv
	case int, int8, int16, int32, int64:
		vv.Raw = vi
		return vv
	case uint, uint8, uint16, uint32, uint64:
		vv.Raw = vi
		return vv
	default:
		return v
	}
}
