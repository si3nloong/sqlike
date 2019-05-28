package options

// LockMode :
type LockMode int

// locks :
const (
	ShareLock LockMode = iota
	UpdateLock
)

// FindOptions :
type FindOptions struct {
	OmitFields []string
	LockMode   LockMode
	IsDebug    bool
}

// SetDebug :
func (opt *FindOptions) SetDebug(debug bool) *FindOptions {
	opt.IsDebug = debug
	return opt
}

// SetOmitFields :
func (opt *FindOptions) SetOmitFields(fields ...string) *FindOptions {
	opt.OmitFields = fields
	return opt
}

// SetLock :
func (opt *FindOptions) SetLock(lm LockMode) *FindOptions {
	opt.LockMode = lm
	return opt
}

// FindOneOptions :
type FindOneOptions = FindOptions
