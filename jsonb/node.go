package jsonb

type jsonType int

const (
	jsonUnknown jsonType = iota
	jsonObject
	jsonArray
	jsonString
	jsonNull
)

// Node :
type Node struct {
	typ     jsonType
	length  int
	prev    *Node
	parent  *Node
	next    *Node
	current *Node
}
