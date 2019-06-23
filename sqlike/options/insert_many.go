package options 

// InsertManyOptions : 
type InsertManyOptions = InsertOptions

// InsertMany :
func InsertMany() *InsertManyOptions {
	return &InsertManyOptions{}
}

// SetMode : 
func (opt *InsertManyOptions) SetMode(mode insertMode) *InsertManyOptions{
	opt.Mode = mode
	return opt
}