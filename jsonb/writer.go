package jsonb

import (
	"bytes"
	"fmt"
	"io"
)

type JsonWriter interface {
	io.StringWriter
	io.ByteWriter
	io.Writer
	fmt.Stringer
	Bytes() []byte
}

// Writer :
type Writer struct {
	bytes.Buffer
}

// NewWriter :
func NewWriter() JsonWriter {
	return &Writer{}
}
