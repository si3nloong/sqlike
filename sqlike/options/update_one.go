package options

// UpdateOneOptions :
type UpdateOneOptions struct {
	Debug bool
}

// UpdateOne :
func UpdateOne() *UpdateOneOptions {
	return &UpdateOneOptions{}
}

// SetDebug :
func (opt *UpdateOneOptions) SetDebug(debug bool) *UpdateOneOptions {
	opt.Debug = debug
	return opt
}
