package primitive

type jsonFunction int

func (f jsonFunction) String() string {
	switch f {
	case JSONContains:
		return "JSON_CONTAINS"
	case JSONExtract:
		return "JSON_EXTRACT"
	case JSONQuote:
		return "JSON_QUOTE"
	case JSONUnquote:
		return "JSON_UNQUOTE"
	case JSONValid:
		return "JSON_VALID"
	case JSONKeys:
		return "JSON_KEYS"
	}
	return "Unknown JSON Function"
}

// sql functions :
const (
	JSONContains jsonFunction = iota + 1
	JSONPretty
	JSONExtract
	JSONQuote
	JSONKeys
	JSONUnquote
	JSONValid
)
