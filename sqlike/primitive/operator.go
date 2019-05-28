package primitive

// Operator :
type Operator int

// operators :
const (
	Equal Operator = iota
	NotEqual
	GreaterThan
	LowerThan
	GreaterEqual
	LowerEqual
	Like
	NotLike
	In
	NotIn
	Between
	NotBetween
	And
	Or
	IsNull
	NotNull
)
