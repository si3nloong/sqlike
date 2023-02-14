package db

import (
	"fmt"
	"io"
)

// Stmt :
type Stmt interface {
	io.StringWriter
	io.ByteWriter
	fmt.Stringer
	Pos() int
	Args() []any
	AppendArgs(args ...any)
	WriteAppendArgs(query string, args ...any)
	Reset()
	StartTimer()
	StopTimer()
}

type Clause interface {
	isSqlClause()
}
