package options

// UpdateManyOptions :
type UpdateManyOptions struct {
	Debug bool
}

// SetDebug :
func (opt *UpdateManyOptions) SetDebug(debug bool) *UpdateManyOptions {
	opt.Debug = debug
	return opt
}
