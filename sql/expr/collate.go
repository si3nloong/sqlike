package expr

import "github.com/si3nloong/sqlike/v2/internal/primitive"

// Collate :
func Collate(collate string, col any, charset ...string) (o primitive.Encoding) {
	if len(charset) > 0 {
		o.Charset = &charset[0]
	}
	o.Column = col
	o.Collate = collate
	return
}
