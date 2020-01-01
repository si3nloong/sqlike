package spatial

import (
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/encoding/wkt"
	"github.com/si3nloong/sqlike/sqlike/primitive"
)

type function int

// functions :
const (
	GeomCollection function = iota + 1
	GeometryCollection
	GeomFromText
	LineString
	MBRContains
	MBRConveredBy
	Distance
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

// // ST_X :
// func ST_X(p interface{}, others ...interface{}) (f Func) {
// 	f.Type = X
// 	return
// }

// // ST_Y :
// func ST_Y(p interface{}, others ...interface{}) (f Func) {
// 	f.Type = Y
// 	return
// }

// // ST_SRID :
// func ST_SRID() (f Func) {
// 	f.Type = SRID
// 	return
// }

// ST_GeomFromText :
func ST_GeomFromText(geo interface{}, sid ...uint) (f Func) {
	f.Type = GeomFromText
	switch vi := geo.(type) {
	case string:
		f.Args = append(f.Args, primitive.Column{
			Name: vi,
		})
	case orb.Geometry:
		f.Args = append(f.Args, primitive.Value{
			Raw: wkt.MarshalString(vi),
		})
	case primitive.Column:
		f.Args = append(f.Args, vi)
	default:
		panic("unsupported data type for ST_GeomFromText")
	}
	if len(sid) > 0 {
		f.Args = append(f.Args, primitive.Value{
			Raw: sid[0],
		})
	}
	return
}

// column, value, ST_GeomFromText(column), ST_GeomFromText(value)
func ST_Distance(g1, g2 interface{}, unit ...string) (f Func) {
	f.Type = Distance
	for _, arg := range []interface{}{g1, g2} {
		switch vi := arg.(type) {
		case string:
		case orb.Geometry:
			f.Args = append(f.Args, primitive.Value{
				Raw: vi,
			})
		case Func:
			f.Args = append(f.Args, vi)
		case primitive.Column:
			f.Args = append(f.Args, vi)
		default:
			panic("unsupported data type for ST_Distance")
		}
	}
	return
}
