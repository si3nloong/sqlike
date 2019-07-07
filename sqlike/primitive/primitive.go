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

// L :
type L string

// Func :
type Func struct {
	Name  string
	Value interface{}
}

// G :
type G []interface{}

// Col :
type Col string

// R :
type R struct {
	From interface{}
	To   interface{}
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

// Sort :
type Sort struct {
	Field string
	Order order
}

// JQ :
type JQ string

// JC :
type JC struct {
	Field interface{}
	Value interface{}
	Path  string
}
