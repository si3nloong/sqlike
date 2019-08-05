package options

// UpdateOptions :
type UpdateOptions struct {
	Debug bool
}

// Update :
func Update() *UpdateOptions {
	return &UpdateOptions{}
}

// SetDebug :
func (opt *UpdateOptions) SetDebug(debug bool) *UpdateOptions {
	opt.Debug = debug
	return opt
}
