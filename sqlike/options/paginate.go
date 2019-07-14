package options

// PaginateOptions :
type PaginateOptions struct {
	Cursor interface{}
	FindOptions
}

// Paginate :
func Paginate() *PaginateOptions {
	return &PaginateOptions{}
}

func (opt *PaginateOptions) SetCursor(cursor interface{}) *PaginateOptions {
	opt.Cursor = cursor
	return opt
}

// SetDebug :
func (opt *PaginateOptions) SetDebug(debug bool) *PaginateOptions {
	opt.Debug = debug
	return opt
}
