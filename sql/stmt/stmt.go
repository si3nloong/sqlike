package sqlstmt

import (
	"sync"
)

var (
	stmtPool = &sync.Pool{
		New: func() interface{} {
			return new(Statement)
		},
	}
)

// AcquireStmt :
func AcquireStmt(fmt Formatter) *Statement {
	x := stmtPool.Get().(*Statement)
	x.fmt = fmt
	return x
}

// ReleaseStmt :
func ReleaseStmt(x *Statement) {
	if x != nil {
		defer stmtPool.Put(x)
		x.Reset()
	}
}
