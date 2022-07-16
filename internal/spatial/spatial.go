package spatial

// Type :
type Type int

// Spatial Type :
const (
	Point Type = iota + 1
	LineString
	Polygon
	MultiPoint
	MultiLineString
	MultiPolygon
)

type function int

// String :
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
	case SpatialTypeTransform:
		return "ST_Transform"
	case SpatialTypeX:
		return "ST_X"
	case SpatialTypeY:
		return "ST_Y"
	case SpatialTypeAsGeoJSON:
		return "ST_AsGeoJSON"
	case SpatialTypeArea:
		return "ST_Area"
	}
	return "UNKNOWN FUNCTION"
}

// functions :
const (
	SpatialTypeGeomFromText function = iota + 1
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
	SpatialTypeAsGeoJSON
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
	Args []any
}

// Geometry :
type Geometry struct {
	Type Type
	SRID uint
	WKT  string
}

// Value :
func (g Geometry) Value() (any, error) {
	return g.WKT, nil
}
