package options

import "github.com/si3nloong/sqlike/sql/util"

type insertMode int

// insert modes :
const (
	InsertIgnore insertMode = iota + 1
	InsertOnDuplicate
)

// InsertOptions :
type InsertOptions struct {
	Mode  insertMode
	Omits util.StringSlice
	Debug bool
}

// Insert :
func Insert() *InsertOptions {
	return &InsertOptions{}
}

// SetMode :
func (opt *InsertOptions) SetMode(mode insertMode) *InsertOptions {
	opt.Mode = mode
	return opt
}

// SetDebug :
func (opt *InsertOptions) SetDebug(debug bool) *InsertOptions {
	opt.Debug = debug
	return opt
}

// SetOmitFields :
func (opt *InsertOptions) SetOmitFields(fields ...string) *InsertOptions {
	opt.Omits = fields
	return opt
}

// // SetOnConflict :
// func (opt *InsertOptions) SetOnConflict(src []interface{}) *InsertOptions {
// 	return opt
// }
