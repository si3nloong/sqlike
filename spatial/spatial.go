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

// Geometry :
type Geometry struct {
	Type Type
	SID  uint
	WKT  string
}

// Value :
func (g Geometry) Value() (interface{}, error) {
	return g.WKT, nil
}
