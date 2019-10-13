package rsql

import "github.com/si3nloong/sqlike/sqlike/primitive"

// Params :
type Params struct {
	Selects []interface{}
	Filters primitive.Group
	Sorts   []interface{}
	Limit   uint
}
