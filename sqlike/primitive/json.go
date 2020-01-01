package primitive

type jsonFunction int

// sql functions :
const (
	JSON_CONTAINS jsonFunction = iota + 1
	JSON_Pretty
	JSON_QUOTE
	JSON_UNQUOTE
	JSON_VALID
)
