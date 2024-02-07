package options

// UpdateOneOptions :
type UpdateOneOptions struct {
	UpdateOptions
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
