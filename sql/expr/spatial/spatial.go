package spatial

type function int

// functions :
const (
	GeomCollection function = iota + 1
	GeometryCollection
	LineString
	MBRContains
	MBRConveredBy
	Point
	Polygon
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

// Func :
type Func struct {
	Type function
	Args []interface{}
}

// ST_X :
func ST_X(p interface{}, others ...interface{}) (f Func) {
	f.Type = X
	return
}

// ST_Y :
func ST_Y(p interface{}, others ...interface{}) (f Func) {
	f.Type = Y
	return
}

// ST_SRID :
func ST_SRID() (f Func) {
	f.Type = SRID
	return
}
