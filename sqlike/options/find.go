package options

// FindOptions :
type FindOptions struct {
	OmitFields []string
	NoLimit    bool
	LockMode   LockMode
	Debug      bool
}

// Find :
func Find() *FindOptions {
	return &FindOptions{}
}

// SetNoLimit :
func (opt *FindOptions) SetNoLimit(limit bool) *FindOptions {
	opt.NoLimit = limit
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

// SetLockMode :
func (opt *FindOptions) SetLockMode(lock LockMode) *FindOptions {
	opt.LockMode = lock
	return opt
}
