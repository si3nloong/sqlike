package options

// UpdateManyOptions :
type UpdateManyOptions struct {
	Debug bool
}

// UpdateMany :
func UpdateMany() *UpdateManyOptions {
	return &UpdateManyOptions{}
}

// SetDebug :
func (opt *UpdateManyOptions) SetDebug(debug bool) *UpdateManyOptions {
	opt.Debug = debug
	return opt
}
