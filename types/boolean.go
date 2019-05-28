package types

import (
	"strconv"
	"strings"
)

// Boolean :
type Boolean bool

// Scan :
func (x *Boolean) Scan(src interface{}) error {
	*x = false
	switch vi := src.(type) {
	case []byte:
		{
			val := strings.ToUpper(string(vi))
			if val == "YES" {
				*x = true
			}
		}

	case string:
		{
			b, _ := strconv.ParseBool(vi)
			*x = Boolean(b)
		}

	case int64:
		{
			if vi == 0 {
				*x = Boolean(false)
			} else if vi == 1 {
				*x = Boolean(true)
			}
		}
	}
	return nil
}
