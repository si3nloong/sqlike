package options

// ModifyOneOptions :
type ModifyOneOptions struct {
	Omits    []string
	Debug    bool
	NoStrict bool
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

// SetStrict :
func (opt *ModifyOneOptions) SetStrict(strict bool) *ModifyOneOptions {
	opt.NoStrict = !strict
	return opt
}

// AppendOmitField :
func (opt *ModifyOneOptions) AppendOmitField(field string) *ModifyOneOptions {
	opt.Omits = append(opt.Omits, field)
	return opt
}
