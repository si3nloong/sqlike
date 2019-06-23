package options

// DeleteManyOptions :
type DeleteManyOptions struct {
	Debug bool
}

// DeleteMany :
func DeleteMany() *DeleteManyOptions {
	return &DeleteManyOptions{}
}

// SetDebug :
func (opt *DeleteManyOptions) SetDebug(debug bool) *DeleteManyOptions {
	opt.Debug = debug
	return opt
}
