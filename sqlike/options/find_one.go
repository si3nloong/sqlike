package options

// // LockMode :
// type LockMode int

// // locks :
// const (
// 	ShareLock LockMode = iota
// 	UpdateLock
// )

// FindOptions :
type FindOptions struct {
	OmitFields []string
	NoLimit    bool
	// LockMode   LockMode
	Debug bool
}

// Find :
func Find() *FindOptions {
	return &FindOptions{}
}

// SetNoLimit :
func (opt *FindOptions) SetNoLimit() *FindOptions {
	opt.NoLimit = true
	return opt
}

// SetDebug :
func (opt *FindOptions) SetDebug(debug bool) *FindOptions {
	opt.Debug = debug
	return opt
}

// SetOmitFields :
func (opt *FindOptions) SetOmitFields(fields ...string) *FindOptions {
	opt.OmitFields = fields
	return opt
}

// SetLock :
// func (opt *FindOptions) SetLock(lm LockMode) *FindOptions {
// 	opt.LockMode = lm
// 	return opt
// }

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
