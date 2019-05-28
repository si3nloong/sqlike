package primitive

// Raw :
type Raw string

// L :
type L string

// G :
type G []interface{}

// Col :
type Col string

// C :
type C struct {
	Field    interface{}
	Operator Operator
	Values   []interface{}
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
