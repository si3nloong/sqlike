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

type function int

func (f function) String() string {
	switch f {
	case ST_GeomFromText:
		return "ST_GeomFromText"
	case ST_Distance:
		return "ST_Distance"
	case ST_Within:
		return "ST_Within"
	case ST_Equals:
		return "ST_Equals"
	case ST_PointFromText:
		return "ST_PointFromText"
	case ST_LineString:
		return "ST_LineString"
	case ST_AsText:
		return "ST_AsText"
	}
	return "UNKNOWN FUNCTION"
}

// functions :
const (
	ST_GeomCollection function = iota + 1
	// GeometryCollection
	ST_GeomFromText
	ST_Distance
	ST_Within
	ST_Equals
	ST_PointFromText
	ST_LineString
	ST_Point
	ST_Polygon
	ST_Area
	ST_AsText
	ST_AsWKB
	ST_AsWKT
	// MBRContains
	// MBRConveredBy
	// AsGeoJSON
	// ST_SRID
	// Transform
	// X
	// Equals
	// Y
	// SymDifference
	// PointFromText
	// PointFromWKB
)

// Func :
type Func struct {
	Type function
	Args []interface{}
}

// Geometry :
type Geometry struct {
	Type Type
	SRID uint
	WKT  string
}

// Value :
func (g Geometry) Value() (interface{}, error) {
	return g.WKT, nil
}
