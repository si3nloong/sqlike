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
	AppendArgs(query string, args ...any)
	Reset()
	StartTimer()
	StopTimer()
}

type SqlClause interface {
	isClause()
}

type SqlStmt interface {
	isStmt()
}
