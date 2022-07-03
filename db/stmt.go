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
}
