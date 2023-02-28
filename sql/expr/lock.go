package expr

import "github.com/si3nloong/sqlike/v2/internal/primitive"

// ForShare :
func ForShare[C ColumnConstraints | *primitive.Lock](opts ...C) (l primitive.Lock) {
	l.Type = primitive.LockForShare
	return
}

// ForUpdate :
func ForUpdate(opts ...func(l *primitive.Lock)) (l primitive.Lock) {
	for _, opt := range opts {
		opt(&l)
	}
	l.Type = primitive.LockForShare
	return
}

func NoWait() func(l *primitive.Lock) {
	return func(l *primitive.Lock) {
		l.Option = primitive.NoWait
	}
}

func SkipLocked() func(l *primitive.Lock) {
	return func(l *primitive.Lock) {
		l.Option = primitive.SkipLocked
	}
}
