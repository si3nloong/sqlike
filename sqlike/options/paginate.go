package options

// PaginateOptions :
type PaginateOptions struct {
	FindOptions
}

// Paginate :
func Paginate() *PaginateOptions {
	return &PaginateOptions{}
}

// SetDebug :
func (opt *PaginateOptions) SetDebug(debug bool) *PaginateOptions {
	opt.Debug = debug
	return opt
}
