package expr

import "github.com/si3nloong/sqlike/sqlike/primitive"

// Collate :
func Collate(collate string, col interface{}, charset ...string) (o primitive.Encoding) {
	if len(charset) > 0 {
		o.Charset = &charset[0]
	}
	o.Column = col
	o.Collate = collate
	return
}

// String :
func String(str string) (o primitive.TypeSafe) {
	o.Type = primitive.Varchar
	o.Value = str
	return
}
