package options

// DeleteManyOptions :
type DeleteManyOptions struct {
	Debug bool
}

// SetDebug :
func (opt *DeleteManyOptions) SetDebug(debug bool) *DeleteManyOptions {
	opt.Debug = debug
	return opt
}
