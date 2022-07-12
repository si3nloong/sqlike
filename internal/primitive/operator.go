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

func (op Operator) String() (n string) {
	switch op {
	case Equal:
		n = "="
	case NotEqual:
		n = "<>"
	case GreaterThan:
		n = ">"
	case LesserThan:
		n = "<"
	case GreaterOrEqual:
		n = ">="
	case LesserOrEqual:
		n = "<="
	case Like:
		n = "LIKE"
	case NotLike:
		n = "NOT LIKE"
	case In:
		n = "IN"
	case NotIn:
		n = "NOT IN"
	case Between:
		n = "BETWEEN"
	case And:
		n = "AND"
	case Or:
		n = "OR"
	case IsNull:
		n = "IS NULL"
	}
	return
}
