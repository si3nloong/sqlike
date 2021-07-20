package expr

import "github.com/si3nloong/sqlike/v2/x/primitive"

// ForShare :
func ForShare(ofs ...primitive.ColumnPath) (l *primitive.Lock) {
	l = new(primitive.Lock)
	l.Type = primitive.LockForShare
	l.Ofs = ofs
	return
}

// ForUpdate :
func ForUpdate(ofs ...primitive.ColumnPath) (l *primitive.Lock) {
	l = new(primitive.Lock)
	l.Type = primitive.LockForUpdate
	l.Ofs = ofs
	return
}
