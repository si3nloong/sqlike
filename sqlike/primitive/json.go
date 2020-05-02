package primitive

type jsonFunction int

// sql functions :
const (
	JSON_CONTAINS jsonFunction = iota + 1
	JSON_PRETTY
	JSON_KEYS
	JSON_TYPE
	JSON_VALID
	JSON_QUOTE
	JSON_UNQUOTE
	JSON_SET
	JSON_EXTRACT
	JSON_INSERT
	JSON_REPLACE
	JSON_REMOVE
)

var jsonFuncNames = [...]string{
	"JSON_CONTAINS",
	"JSON_PRETTY",
	"JSON_KEYS",
	"JSON_TYPE",
	"JSON_VALID",
	"JSON_QUOTE",
	"JSON_UNQUOTE",
	"JSON_SET",
	"JSON_EXTRACT",
	"JSON_INSERT",
	"JSON_REPLACE",
	"JSON_REMOVE",
}

func (f jsonFunction) String() string {
	id := int(f) - 1
	if id > len(jsonFuncNames) {
		return "Unknown JSON Function"
	}
	return jsonFuncNames[id]
}
