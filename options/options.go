package options

import (
	"github.com/si3nloong/sqlike/v2/internal/primitive"
)

// LockForShare :
func LockForShare(opts ...func(l *primitive.Lock)) (l primitive.Lock) {
	for _, opt := range opts {
		opt(&l)
	}
	l.Type = primitive.LockForShare
	return
}

// LockForUpdate :
func LockForUpdate(opts ...func(l *primitive.Lock)) (l primitive.Lock) {
	for _, opt := range opts {
		opt(&l)
	}
	l.Type = primitive.LockForUpdate
	return
}
