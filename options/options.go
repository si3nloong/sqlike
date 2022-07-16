package options

import (
	"github.com/si3nloong/sqlike/v2/internal/primitive"
)

// LockForUpdate :
func LockForUpdate() (l primitive.Lock) {
	l.Type = primitive.LockForUpdate
	return
}

// LockForShare :
func LockForShare() (l primitive.Lock) {
	l.Type = primitive.LockForShare
	return
}
