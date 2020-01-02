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
	case SpatialTypeGeomFromText:
		return "ST_GeomFromText"
	case SpatialTypeDistance:
		return "ST_Distance"
	case SpatialTypeWithin:
		return "ST_Within"
	case SpatialTypeEquals:
		return "ST_Equals"
	case SpatialTypePointFromText:
		return "ST_PointFromText"
	case SpatialTypeLineString:
		return "ST_LineString"
	case SpatialTypeAsText:
		return "ST_AsText"
	}
	return "UNKNOWN FUNCTION"
}

// functions :
const (
	ST_GeomCollection function = iota + 1
	// GeometryCollection
	SpatialTypeGeomFromText
	SpatialTypeDistance
	SpatialTypeWithin
	SpatialTypeEquals
	SpatialTypePointFromText
	SpatialTypeLineString
	SpatialTypePoint
	SpatialTypePolygon
	SpatialTypeArea
	SpatialTypeAsText
	SpatialTypeAsWKB
	SpatialTypeAsWKT
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
