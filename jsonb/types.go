package jsonb

type jsonType int

const (
	jsonInvalid jsonType = iota
	jsonNull
	jsonObject
	jsonArray
	jsonWhitespace
	jsonLiteral
	jsonString
	jsonComma
	jsonBoolean
	jsonNumber
)

func (jt jsonType) String() (name string) {
	switch jt {
	case jsonInvalid:
		name = "invalid"
	case jsonNull:
		name = "null"
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
