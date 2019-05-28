package indexes

// Kind :
type Kind int

// types :
const (
	FullText Kind = iota + 1
	Unique
	Spatial
)

// Order :
type Order int

// ordering
const (
	Desc Order = iota
	Asc
)

// Index :
type Index struct {
	Name    string
	Kind    Kind
	Order   Order
	Columns []string
}
