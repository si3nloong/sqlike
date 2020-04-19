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

func (t Type) String() string {
	switch t {
	case String:
		return "string"
	case Bool:
		return "boolean"
	case Byte:
		return "byte"
	case Int:
		return "int"
	case Int8:
		return "int8"
	case Int16:
		return "int16"
	case Int32:
		return "int32"
	case Int64:
		return "int64"
	case Uint:
		return "uint"
	case Uint8:
		return "uint8"
	case Uint16:
		return "uint16"
	case Uint32:
		return "uint32"
	case Uint64:
		return "uint64"
	case Float32:
		return "float32"
	case Float64:
		return "float64"
	case Slice:
		return "slice"
	case Map:
		return "map"
	case Struct:
		return "struct"
	case Timestamp:
		return "timestamp"
	case DateTime:
		return "datetime"
	case JSON:
		return "json"
	case UUID:
		return "uuid"
	case Point:
		return "point"
	case LineString:
		return "linestring"
	case Polygon:
		return "polygon"
	case MultiPoint:
		return "multipoint"
	case MultiLineString:
		return "multilinestring"
	case MultiPolygon:
		return "multipolygon"
	default:
		return "unknown"
	}
}
