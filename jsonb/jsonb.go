package jsonb

var (
	valueMap = make([]jsonType, 256)
)

func init() {
	length := len(valueMap)
	for i := 0; i < length; i++ {
		valueMap[i] = jsonInvalid
	}
	valueMap['"'] = jsonString
	valueMap['-'] = jsonNumber
	valueMap['0'] = jsonNumber
	valueMap['1'] = jsonNumber
	valueMap['2'] = jsonNumber
	valueMap['3'] = jsonNumber
	valueMap['4'] = jsonNumber
	valueMap['5'] = jsonNumber
	valueMap['6'] = jsonNumber
	valueMap['7'] = jsonNumber
	valueMap['8'] = jsonNumber
	valueMap['9'] = jsonNumber
	valueMap['t'] = jsonBoolean
	valueMap['f'] = jsonBoolean
	valueMap['n'] = jsonNull
	valueMap['['] = jsonArray
	valueMap['{'] = jsonObject
	valueMap[' '] = jsonWhitespace
	valueMap['\r'] = jsonWhitespace
	valueMap['\t'] = jsonWhitespace
	valueMap['\n'] = jsonWhitespace
}
