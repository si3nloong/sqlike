package spatial

// Type
type Type int

const (
	Point Type = iota + 1
	LineString
	Polygon
	MultiPoint
	MultiLineString
	MultiPolygon
)

// Function :
type Function int

// functions :
const (
	GeomCollection Function = iota + 1
	GeometryCollection
	ST_LineString
	MBRContains
	MBRConveredBy
	ST_Point
	ST_Polygon
	Area
	AsText
	AsWKB
	AsWKT
	AsGeoJSON
	SRID
	Transform
	Equals
	Within
	X
	Y
	SymDifference
	PointFromText
	PointFromWKB
)

// Geometry :
type Geometry struct {
	Type  Type
	Value string
}
