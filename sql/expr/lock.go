package expr

import "github.com/si3nloong/sqlike/v2/internal/primitive"

// ForShare :
func ForShare[T primitive.ColumnPath | string](ofs ...T) (l *primitive.Lock) {
	l = new(primitive.Lock)
	l.Type = primitive.LockForShare
	// l.Ofs = ofs
	return
}

// ForUpdate :
func ForUpdate[T primitive.ColumnPath | string](ofs ...T) (l *primitive.Lock) {
	l = new(primitive.Lock)
	l.Type = primitive.LockForUpdate
	// l.Ofs = ofs
	return
}
