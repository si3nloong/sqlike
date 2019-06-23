package options

// UpdateOneOptions :
type UpdateOneOptions struct {
	Debug bool
}

// SetDebug :
func (opt *UpdateOneOptions) SetDebug(debug bool) *UpdateOneOptions {
	opt.Debug = debug
	return opt
}
