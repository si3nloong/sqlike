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
	case SpatialTypeAsWKB:
		return "ST_AsWKB"
	case SpatialTypeAsWKT:
		return "ST_AsWKT"
	case SpatialTypeSRID:
		return "ST_SRID"
	case SpatialTypeIsValid:
		return "ST_IsValid"
	case SpatialTypeIntersects:
		return "ST_Intersects"
	}
	return "UNKNOWN FUNCTION"
}

// functions :
const (
	SpatialTypeGeomCollection function = iota + 1
	SpatialTypeGeomFromText
	SpatialTypeDistance
	SpatialTypeWithin
	SpatialTypeEquals
	SpatialTypePointFromText
	SpatialTypePointFromWKB
	SpatialTypeLineString
	SpatialTypePoint
	SpatialTypePolygon
	SpatialTypeArea
	SpatialTypeAsText
	SpatialTypeAsWKB
	SpatialTypeAsWKT
	SpatilaTypeAsGeoJSON
	SpatialTypeSRID
	SpatialTypeX
	SpatialTypeY
	SpatialTypeIsValid
	SpatialTypeIntersects
	SpatialTypeTransform
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
