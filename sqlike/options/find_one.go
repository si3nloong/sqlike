package options

// FindOneOptions :
type FindOneOptions struct {
	FindOptions
}

// FindOne :
func FindOne() *FindOneOptions {
	return &FindOneOptions{}
}

// SetDebug :
func (opt *FindOneOptions) SetDebug(debug bool) *FindOneOptions {
	opt.Debug = debug
	return opt
}

// SetOmitFields :
func (opt *FindOneOptions) SetOmitFields(fields ...string) *FindOneOptions {
	opt.OmitFields = fields
	return opt
}
