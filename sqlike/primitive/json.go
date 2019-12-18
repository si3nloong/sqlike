package primitive

type jsonFunction int

// sql functions :
const (
	JSONArray jsonFunction = iota + 1
	JSONContains
	JSONDepth
	JSONExtract
	JSONInsert
	JSONKeys
	JSONLength
	JSONObject
	JSONPretty
	JSONQuote
	JSONRemove
	JSONReplace
	JSONSearch
	JSONSet
	JSONType
	JSONUnquote
	JSONValid
)
