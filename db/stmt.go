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
	Args() []any
	AppendArgs(args ...any)
	WriteAppendArgs(query string, args ...any)
	Reset()
	StartTimer()
	StopTimer()
}

// this is to enable struct to implement sql clause interface
type BaseClause struct{}

func (b BaseClause) isSqlClause() {}

type Clause interface {
	isSqlClause()
}
