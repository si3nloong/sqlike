package options 

// InsertManyOptions : 
type InsertManyOptions = InsertOptions

// SetDebug : 
func (opt *InsertManyOptions) SetDebug(debug bool) *InsertManyOptions{
	opt.IsDebug = debug
	return opt
}

// SetMode : 
func (opt *InsertManyOptions) SetMode(mode insertMode) *InsertManyOptions{
	opt.Mode = mode
	return opt
}