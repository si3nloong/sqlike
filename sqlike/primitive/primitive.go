package primitive

// Raw :
type Raw struct {
	Value string
}

// Column :
type Column struct {
	Table string
	Name  string
}

// Alias :
type Alias struct {
	Name  string
	Alias string
}

// CastAs :
type CastAs struct {
	Value    interface{}
	DataType DataType
}

// Func :
type Func struct {
	Type      Function
	Arguments []interface{}
}

// Group :
type Group []interface{}

// Col :
type Col string

// R :
type R struct {
	From interface{}
	To   interface{}
}

// Field :
type Field struct {
	Name   string
	Values []interface{}
}

// L :
type L struct {
	Field interface{}
	IsNot bool
	Value interface{}
}

// C :
type C struct {
	Field    interface{}
	Operator Operator
	Value    interface{}
}

// KV :
type KV struct {
	Field string
	Value interface{}
}

type operator int

// operators :
const (
	Add operator = iota
	Deduct
)

// Math :
type Math struct {
	Field Col
	Mode  operator
	Value int
}

type order int

// orders :
const (
	Ascending order = iota
	Descending
)

// Nil :
type Nil struct {
	Field interface{}
	IsNot bool
}

// Sort :
type Sort struct {
	Field string
	Order order
}

// Value :
type Value struct {
	Raw interface{}
}

// JC :
type JC struct {
	Target    interface{}
	Candidate interface{}
	Path      *string
}

type aggregate int

// aggregation :
const (
	Sum aggregate = iota + 1
	Count
	Average
	Max
	Min
)

// Aggregate :
type Aggregate struct {
	Field interface{}
	By    aggregate
}

// As :
type As struct {
	Field interface{}
	Name  string
}
