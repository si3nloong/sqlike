package examples

import (
	"database/sql"
	"testing"

	"github.com/paulmach/orb"
	"github.com/si3nloong/sqlike/sql/expr"
	"github.com/si3nloong/sqlike/sqlike"
	"github.com/si3nloong/sqlike/sqlike/actions"
	"github.com/si3nloong/sqlike/sqlike/options"
	"github.com/stretchr/testify/require"
)

type Spatial struct {
	ID    int64 `sqlike:",primary_key"`
	Point orb.Point
	// LineString orb.LineString
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
		sp.Point = point
		// sp.LineString = []orb.Point{
		// 	orb.Point{0, 0},
		// 	orb.Point{1, 1},
		// 	orb.Point{2, 2},
		// }
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
		_, err = table.InsertOne(&sp,
			options.InsertOne().SetDebug(true),
		)
		require.NoError(t, err)
	}

	var o struct {
		ID    int64
		Point orb.Point
	}

	{

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
	}

	{
		origin := orb.Point{1, 5}
		err = table.FindOne(
			actions.FindOne().
				Select(
					expr.Func("ST_Distance", expr.Column("Point"), origin),
				).
				Where(
					expr.Equal("ID", 1),
				),
			options.FindOne().SetDebug(true),
		).Decode(&o)
		require.NoError(t, err)
		// values := make([]interface{}, 2)
		// result.Scan(values...)
	}
}
