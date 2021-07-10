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
	Args() []interface{}
	AppendArgs(args ...interface{})
	WriteAppendArgs(query string, args ...interface{})
}
