package options

// DestroyOneOptions :
type DestroyOneOptions struct {
	DeleteOptions
}

// DestroyOne :
func DestroyOne() *DestroyOneOptions {
	return &DestroyOneOptions{}
}

// SetDebug :
func (opt *DestroyOneOptions) SetDebug(debug bool) *DestroyOneOptions {
	opt.Debug = debug
	return opt
}
