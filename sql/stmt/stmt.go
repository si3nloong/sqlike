package sqlstmt

import (
	"sync"
)

var (
	stmtPool = &sync.Pool{
		New: func() any {
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
		// this will reset everything including timer, query statement and arguments
		x.Reset()
		x.fmt = nil
		defer stmtPool.Put(x)
	}
}
