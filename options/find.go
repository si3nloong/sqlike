package options

import "github.com/si3nloong/sqlike/v2/internal/primitive"

// FindOptions :
type FindOptions struct {
	OmitFields []string
	NoLimit    bool
	Lock       primitive.Lock
	Debug      bool
}

// Find :
func Find() *FindOptions {
	return &FindOptions{Lock: primitive.Lock{}}
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
func (opt *FindOptions) SetLockMode(lock primitive.Lock) *FindOptions {
	opt.Lock = lock
	return opt
}
