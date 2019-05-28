package options

// ModifyOneOptions :
type ModifyOneOptions struct {
	IsDebug bool
}

// ModifyOne :
func ModifyOne() *ModifyOneOptions {
	return &ModifyOneOptions{}
}

// SetDebug :
func (opt *ModifyOneOptions) SetDebug(debug bool) *ModifyOneOptions {
	opt.IsDebug = debug
	return opt
}
