package jsonb

import "bytes"

type jsonType int

const (
	object jsonType = iota
	array
)

// Reader :
type Reader struct {
	bytes.Buffer
}

// Node :
type Node struct {
	typ     jsonType
	length  int
	prev    *Node
	parent  *Node
	next    *Node
	current *Node
}
