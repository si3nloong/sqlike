package jsonb

type jsonType int

// json type values :
const (
	jsonInvalid jsonType = iota
	jsonNull
	jsonObject
	jsonArray
	jsonWhitespace
	jsonString
	jsonBoolean
	jsonNumber
	// jsonLiteral
	// jsonComma
)

func (jt jsonType) String() (name string) {
	switch jt {
	case jsonInvalid:
		name = "invalid"
	case jsonNull:
		name = "null"
	case jsonWhitespace:
		name = "whitespace"
	case jsonString:
		name = "string"
	case jsonBoolean:
		name = "boolean"
	case jsonObject:
		name = "object"
	case jsonArray:
		name = "array"
	case jsonNumber:
		name = "number"
	}
	return
}
