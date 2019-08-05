package options

// DeleteOneOptions :
type DeleteOneOptions struct {
	DeleteManyOptions
}

// DeleteOne :
func DeleteOne() *DeleteOneOptions {
	return &DeleteOneOptions{}
}

// SetDebug :
func (opt *DeleteOneOptions) SetDebug(debug bool) *DeleteOneOptions {
	opt.Debug = debug
	return opt
}
