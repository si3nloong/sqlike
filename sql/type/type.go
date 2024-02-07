package sqltype

// Type : merge the golang data type and custom type
type Type int

// types :
const (
	String Type = iota
	Char
	Bool
	Byte
	Int
	Int8
	Int16
	Int32
	Int64
	Uint
	Uint8
	Uint16
	Uint32
	Uint64
	Float32
	Float64
	GeoPoint
	Date
	DateTime
	Time
	Timestamp
	Struct
	Array
	Slice
	Map
	JSON
	Enum
	Set
	UUID
	Point
	LineString
	Polygon
	MultiPoint
	MultiLineString
	MultiPolygon
)

var names = map[Type]string{
	String:          "string",
	Bool:            "boolean",
	Byte:            "byte",
	Int:             "int",
	Int8:            "int8",
	Int16:           "int16",
	Int32:           "int32",
	Int64:           "int64",
	Uint:            "uint",
	Uint8:           "uint8",
	Uint16:          "uint16",
	Uint32:          "uint32",
	Uint64:          "uint64",
	Float32:         "float32",
	Float64:         "float64",
	Slice:           "slice",
	Map:             "map",
	Struct:          "struct",
	Timestamp:       "timestamp",
	DateTime:        "datetime",
	Time:            "time",
	JSON:            "json",
	UUID:            "uuid",
	Point:           "point",
	LineString:      "linestring",
	Polygon:         "polygon",
	MultiPoint:      "multipoint",
	MultiLineString: "multilinestring",
	MultiPolygon:    "multipolygon",
}

// String :
func (t Type) String() string {
	if v, ok := names[t]; ok {
		return v
	}
	return "unknown type"
}
