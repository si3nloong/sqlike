package options

// ModifyOneOptions :
type ModifyOneOptions struct {
	Omits []string
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

// SetOmitFields :
func (opt *ModifyOneOptions) SetOmitFields(fields ...string) *ModifyOneOptions {
	opt.Omits = fields
	return opt
}
