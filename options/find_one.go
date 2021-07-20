package options

import "github.com/si3nloong/sqlike/v2/x/primitive"

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

// SetLockMode :
func (opt *FindOneOptions) SetLockMode(lock primitive.Lock) *FindOneOptions {
	opt.Lock = lock
	return opt
}
