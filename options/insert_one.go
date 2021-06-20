package options

// InsertOneOptions :
type InsertOneOptions struct {
	InsertOptions
}

// InsertOne :
func InsertOne() *InsertOneOptions {
	return &InsertOneOptions{}
}

// SetMode :
func (opt *InsertOneOptions) SetMode(mode insertMode) *InsertOneOptions {
	opt.Mode = mode
	return opt
}

// SetDebug :
func (opt *InsertOneOptions) SetDebug(debug bool) *InsertOneOptions {
	opt.Debug = debug
	return opt
}

// SetOmitFields :
func (opt *InsertOneOptions) SetOmitFields(fields ...string) *InsertOneOptions {
	opt.Omits = fields
	return opt
}

// // SetOnConflict :
// func (opt *InsertOneOptions) SetOnConflict(src []interface{}) *InsertOneOptions {
// 	return opt
// }
