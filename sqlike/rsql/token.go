package rsql

type tokenType int

// types :
const (
	Operator = iota
	String
	Or
	And
	Numeric
	Text
	Group
	Whitespace
)
