package expr

import (
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/encoding/wkt"
	"github.com/si3nloong/sqlike/spatial"
	"github.com/si3nloong/sqlike/sqlike/primitive"
)

// ST_GeomFromText :
func ST_GeomFromText(g interface{}, srid ...uint) (f spatial.Func) {
	f.Type = spatial.ST_GeomFromText
	switch vi := g.(type) {
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
	if len(srid) > 0 {
		f.Args = append(f.Args, primitive.Value{
			Raw: srid[0],
		})
	}
	return
}

// ST_AsText :
func ST_AsText(g interface{}) (f spatial.Func) {
	f.Type = spatial.ST_AsText
	switch vi := g.(type) {
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
		panic("unsupported data type for ST_AsText")
	}
	return
}

// column, value, ST_GeomFromText(column), ST_GeomFromText(value)
// ST_Distance :
func ST_Distance(g1, g2 interface{}, unit ...string) (f spatial.Func) {
	f.Type = spatial.ST_Distance
	for _, arg := range []interface{}{g1, g2} {
		switch vi := arg.(type) {
		case string:
		case orb.Geometry:
			f.Args = append(f.Args, primitive.Value{
				Raw: vi,
			})
		case spatial.Func:
			f.Args = append(f.Args, vi)
		case primitive.Column:
			f.Args = append(f.Args, vi)
		default:
			panic("unsupported data type for ST_Distance")
		}
	}
	return
}

// ST_Equals :
func ST_Equals(g1, g2 interface{}) (f spatial.Func) {
	f.Type = spatial.ST_Equals
	for _, arg := range []interface{}{g1, g2} {
		switch vi := arg.(type) {
		case string:
		case orb.Geometry:
			f.Args = append(f.Args, primitive.Value{
				Raw: vi,
			})
		case spatial.Func:
			f.Args = append(f.Args, vi)
		case primitive.Column:
			f.Args = append(f.Args, vi)
		default:
			panic("unsupported data type for ST_Equals")
		}
	}
	return
}

// ST_Within :
func ST_Within(g1, g2 interface{}) (f spatial.Func) {
	f.Type = spatial.ST_Within
	for _, arg := range []interface{}{g1, g2} {
		switch vi := arg.(type) {
		case string:
		case orb.Geometry:
			f.Args = append(f.Args, primitive.Value{
				Raw: vi,
			})
		case spatial.Func:
			f.Args = append(f.Args, vi)
		case primitive.Column:
			f.Args = append(f.Args, vi)
		default:
			panic("unsupported data type for ST_Within")
		}
	}
	return
}
