package types

import (
	"strconv"
	"strings"
)

// Boolean :
type Boolean bool

// Scan :
func (x *Boolean) Scan(src any) error {
	*x = false
	switch vi := src.(type) {
	case []byte:
		{
			val := strings.ToLower(string(vi))
			switch val {
			case "yes", "y":
				*x = true
			case "no", "n":
				*x = false
			default:
				b, _ := strconv.ParseBool(val)
				*x = Boolean(b)
			}
		}

	case string:
		{
			vi = strings.ToLower(vi)
			switch vi {
			case "yes", "y":
				*x = true
			case "no", "n":
				*x = false
			default:
				b, _ := strconv.ParseBool(vi)
				*x = Boolean(b)
			}
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
