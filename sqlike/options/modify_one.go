package options

// ModifyOneOptions :
type ModifyOneOptions struct {
	Debug bool
}

// ModifyOne :
func ModifyOne() *ModifyOneOptions {
	return &ModifyOneOptions{}
}

// SetDebug :
func (opt *ModifyOneOptions) SetDebug(debug bool) *ModifyOneOptions {
	opt.Debug = debug
	return opt
}
