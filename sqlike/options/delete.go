package options

// DeleteOptions :
type DeleteOptions struct {
	Debug bool
}

// Delete :
func Delete() *DeleteOptions {
	return &DeleteOptions{}
}

// SetDebug :
func (opt *DeleteOptions) SetDebug(debug bool) *DeleteOptions {
	opt.Debug = debug
	return opt
}
