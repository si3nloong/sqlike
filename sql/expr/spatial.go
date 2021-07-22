package expr

import (
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/encoding/wkt"
	"github.com/si3nloong/sqlike/v2/x/primitive"
	"github.com/si3nloong/sqlike/v2/x/spatial"
)

// ST_GeomFromText :
func ST_GeomFromText(g interface{}, srid ...uint) (f spatial.Func) {
	f.Type = spatial.SpatialTypeGeomFromText
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
	f.Type = spatial.SpatialTypeAsText
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

// ST_IsValid :
func ST_IsValid(g interface{}) (f spatial.Func) {
	f.Type = spatial.SpatialTypeIsValid
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
		panic("unsupported data type for ST_IsValid")
	}
	return
}

// column, value, ST_GeomFromText(column), ST_GeomFromText(value)
// ST_Distance :
func ST_Distance(g1, g2 interface{}, unit ...string) (f spatial.Func) {
	f.Type = spatial.SpatialTypeDistance
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
	f.Type = spatial.SpatialTypeEquals
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

// ST_Intersects :
func ST_Intersects(g1, g2 interface{}) (f spatial.Func) {
	f.Type = spatial.SpatialTypeIntersects
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
			panic("unsupported data type for ST_Intersects")
		}
	}
	return
}

// ST_Within :
func ST_Within(g1, g2 interface{}) (f spatial.Func) {
	f.Type = spatial.SpatialTypeWithin
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
