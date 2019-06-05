package jsonb

import "bytes"

// Writer :
type Writer struct {
	bytes.Buffer
}

// NewWriter :
func NewWriter() *Writer {
	return &Writer{}
}
