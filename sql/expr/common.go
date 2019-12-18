package expr

import (
	"database/sql"
	"strings"

	"github.com/si3nloong/sqlike/sqlike/primitive"
)

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

// Function :
func Function(name string, value interface{}, others ...interface{}) (f primitive.Func) {
	f.Name = strings.ToUpper(strings.TrimSpace(name))
	f.Arguments = append(f.Arguments, wrapRaw(value))
	if len(others) > 0 {
		for _, arg := range others {
			f.Arguments = append(f.Arguments, wrapRaw(arg))
		}
	}
	return
}

func wrapRaw(v interface{}) (it interface{}) {
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
