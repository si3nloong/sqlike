package options 

// InsertOneOptions : 
type InsertOneOptions = InsertOptions

// SetOnConflict : 
func (opt *InsertOneOptions) SetOnConflict(src []interface{}) *InsertOneOptions {
	return opt
}

// AppendOnConflict : 
func (opt *InsertOneOptions) AppendOnConflict(col string, val interface{}) *InsertOneOptions {
	return opt
}