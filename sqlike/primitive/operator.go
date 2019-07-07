package primitive

// Operator :
type Operator int

// operators :
const (
	Equal Operator = iota
	NotEqual
	GreaterThan
	LesserThan
	GreaterOrEqual
	LesserOrEqual
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
