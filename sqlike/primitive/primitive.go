package primitive

// Raw :
type Raw string

// L :
type L string

// G :
type G []interface{}

// GV :
type GV []interface{}

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
	Field Col
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
