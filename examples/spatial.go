package examples

import (
	"database/sql"
	"testing"

	"github.com/paulmach/orb"
	"github.com/si3nloong/sqlike/sql/expr"
	"github.com/si3nloong/sqlike/sql/expr/spatial"
	"github.com/si3nloong/sqlike/sqlike"
	"github.com/si3nloong/sqlike/sqlike/actions"
	"github.com/si3nloong/sqlike/sqlike/options"
	"github.com/stretchr/testify/require"
)

type Spatial struct {
	ID             int64 `sqlike:",primary_key"`
	Point          orb.Point
	PtrPoint       *orb.Point
	Point4326      orb.Point `sqlike:"PointWithSID,sid=4326"`
	LineString     orb.LineString
	LineString2    orb.LineString
	LineString3    orb.LineString
	PtrLineString  *orb.LineString
	LineString4326 orb.LineString `sqlike:"LineStringWithSID,sid=4326"`
	// Polygon    orb.Polygon
}

// SpatialExamples :
func SpatialExamples(t *testing.T, db *sqlike.Database) {
	var (
		sp    = Spatial{}
		table = db.Table("spatial")
		err   error
	)

	point := orb.Point{1, 5}

	{
		err = table.DropIfExits()
		require.NoError(t, err)
	}

	{
		table.MustMigrate(Spatial{})
	}

	{
		sp.ID = 1
		sp.Point4326 = point
		sp.Point = point
		sp.LineString = []orb.Point{
			orb.Point{0, 0},
			orb.Point{1, 1},
		}
		sp.LineString2 = []orb.Point{
			orb.Point{0, 0},
			orb.Point{1, 1},
			orb.Point{2, 2},
		}
		sp.LineString3 = []orb.Point{
			orb.Point{0, 0},
			orb.Point{1, 1},
			orb.Point{2, 2},
			orb.Point{3, 3},
			orb.Point{4, 4},
		}
		sp.LineString4326 = []orb.Point{
			orb.Point{88, 0},
			orb.Point{1, 10},
		}
		// sp.Polygon = orb.Polygon{
		// 	// (0 0,10 0,10 10,0 10,0 0)
		// 	orb.Ring{
		// 		orb.Point{0, 0},
		// 		orb.Point{10, 0},
		// 		orb.Point{10, 10},
		// 		orb.Point{0, 10},
		// 		orb.Point{0, 0},
		// 	},
		// 	// (5 5,7 5,7 7,5 7, 5 5)
		// 	orb.Ring{
		// 		orb.Point{5, 5},
		// 		orb.Point{7, 5},
		// 		orb.Point{7, 7},
		// 		orb.Point{5, 7},
		// 		orb.Point{5, 5},
		// 	},
		// }
		sps := []Spatial{sp, sp, sp}
		_, err = table.Insert(&sps,
			options.Insert().SetDebug(true),
		)
		require.NoError(t, err)
	}

	{
		var o Spatial
		result := table.FindOne(
			actions.FindOne().Where(
				expr.Equal("ID", 1),
			),
			options.FindOne().SetDebug(true),
		)
		b := new(sql.RawBytes)
		err = result.Scan(b)
		require.NoError(t, err)
		err = result.Decode(&o)
		require.NoError(t, err)

		require.Equal(t, int64(1), o.ID)
		require.Equal(t, point, o.Point)
		require.Equal(t, orb.Point{5, 1}, o.Point4326)
	}

	{
		origin := orb.Point{20, 10}
		p1 := orb.Point{1, 3}
		p2 := orb.Point{4, 18}
		var dist1, dist2 float64
		err = table.FindOne(
			actions.FindOne().
				Select(
					expr.As(spatial.ST_Distance(expr.Column("Point"), origin), "dist"),
					spatial.ST_Distance(
						spatial.ST_GeomFromText(p1, 4326),
						spatial.ST_GeomFromText(p2, 4326),
					),
				).
				Where(
					expr.Equal("ID", 1),
				).
				OrderBy(
					expr.Desc("dist"),
				),
			options.FindOne().SetDebug(true),
		).Scan(&dist1, &dist2)
		require.NoError(t, err)
		require.Equal(t, float64(19.6468827043885), dist1)
		require.True(t, dist2 > 0)
	}
}
