package expr

import (
	"github.com/si3nloong/sqlike/sqlike/primitive"
)

// JSONQuote :
func JSONQuote(value string) primitive.JQ {
	return primitive.JQ(value)
}

// JSONContains :
func JSONContains(field, value interface{}, paths ...string) (jc primitive.JC) {
	path := "$"
	if len(paths) > 0 {
		path = paths[0]
	}
	jc.Field = field
	jc.Value = value
	jc.Path = path
	return
}
