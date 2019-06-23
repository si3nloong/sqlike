package options 

// InsertOneOptions : 
type InsertOneOptions = InsertOptions

// InsertOne :
func InsertOne() *InsertOneOptions {
	return &InsertOneOptions{}
}

// SetDebug : 
func (opt *InsertOneOptions) SetDebug(debug bool) *InsertOneOptions{
	opt.Debug = debug
	return opt
}

// SetOnConflict : 
func (opt *InsertOneOptions) SetOnConflict(src []interface{}) *InsertOneOptions {
	return opt
}

// AppendOnConflict : 
func (opt *InsertOneOptions) AppendOnConflict(col string, val interface{}) *InsertOneOptions {
	return opt
}