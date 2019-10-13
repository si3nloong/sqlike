package rsql

import (
	"net/url"
	"strings"
)

func parseRawQuery(m map[string]string, query string) (err error) {
	for query != "" {
		key := query
		if i := strings.IndexAny(key, "&"); i >= 0 {
			key, query = key[:i], key[i+1:]
		} else {
			query = ""
		}
		if key == "" {
			continue
		}
		value := ""
		if i := strings.Index(key, "="); i >= 0 {
			key, value = key[:i], key[i+1:]
		}
		key, err1 := url.QueryUnescape(key)
		if err1 != nil {
			if err == nil {
				err = err1
			}
			continue
		}
		value, err1 = url.QueryUnescape(value)
		if err1 != nil {
			if err == nil {
				err = err1
			}
			continue
		}
		// m[key] = append(m[key], value)
		m[key] = value
	}
	return err
}

// // expression  = [ "(" ]
// // ( constraint / expression )
// // [ operator ( constraint / expression ) ]
// // [ ")" ]
// // operator    = ";" / ","
