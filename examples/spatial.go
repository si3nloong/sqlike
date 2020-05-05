package examples

import (
	"context"
	"database/sql"
	"testing"

	"github.com/paulmach/orb"
	"github.com/si3nloong/sqlike/sql/expr"
	"github.com/si3nloong/sqlike/sqlike"
	"github.com/si3nloong/sqlike/sqlike/actions"
	"github.com/si3nloong/sqlike/sqlike/indexes"
	"github.com/si3nloong/sqlike/sqlike/options"
	"github.com/stretchr/testify/require"
)

type Spatial struct {
	ID             int64 `sqlike:",primary_key,auto_increment"`
	Point          orb.Point
	PtrPoint       *orb.Point
	Point4326      orb.Point `sqlike:"PointWithSRID,srid=4326"`
	LineString     orb.LineString
	LineString2    orb.LineString
	LineString3    orb.LineString
	PtrLineString  *orb.LineString
	LineString4326 orb.LineString `sqlike:"LineStringWithSRID,srid=4326"`
	// Polygon    orb.Polygon
}

// SpatialExamples :
func SpatialExamples(t *testing.T, ctx context.Context, db *sqlike.Database) {
	var (
		sp    = Spatial{}
		table = db.Table("spatial")
		err   error
	)

	point := orb.Point{1, 5}

	{
		err = table.DropIfExists(ctx)
		require.NoError(t, err)
	}

	{
		table.MustMigrate(ctx, Spatial{})
		table.MustUnsafeMigrate(ctx, Spatial{})
		iv := table.Indexes()
		idx := indexes.Index{
			Type:    indexes.Spatial,
			Columns: indexes.Columns("Point"),
		}
		err = iv.CreateOne(ctx, idx)
		require.NoError(t, err)
		result, err := iv.List(ctx)
		require.NoError(t, err)
		require.True(t, len(result) > 0)
		require.Equal(t, sqlike.Index{
			Name:     idx.GetName(),
			Type:     "SPATIAL",
			IsUnique: false,
		}, result[0])
	}

	{
		sp.Point4326 = point
		sp.Point = point
		sp.LineString = []orb.Point{
			{0, 0},
			{1, 1},
		}
		sp.LineString2 = []orb.Point{
			{0, 0},
			{1, 1},
			{2, 2},
		}
		sp.LineString3 = []orb.Point{
			{0, 0},
			{1, 1},
			{2, 2},
			{3, 3},
			{4, 4},
		}
		sp.LineString4326 = []orb.Point{
			{88, 0},
			{1, 10},
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
		_, err = table.Insert(
			ctx,
			&sps,
			options.Insert().SetDebug(true),
		)
		require.NoError(t, err)
	}

	{
		result := table.FindOne(
			ctx,
			actions.FindOne().
				Select(
					"ID",
					"Point",
					"PtrPoint",
				).
				Where(
					expr.Equal("ID", 1),
				),
			options.FindOne().SetDebug(true),
		)
		// b := new(sql.RawBytes)
		// var str string

		c1 := new(sql.RawBytes)
		c2 := new(sql.RawBytes)
		c3 := new(sql.RawBytes)
		cols := result.Columns()
		require.ElementsMatch(t, []string{
			"ID",
			"Point",
			"PtrPoint",
		}, cols)
		err = result.Scan(c1, c2, c3)
		require.NoError(t, err)

		v1 := sql.RawBytes(`1`)
		require.Equal(t, &v1, c1)
		// TODO: check column2 value
		require.Equal(t, &v1, c1)
		require.Nil(t, *c3)
	}

	{
		var o Spatial
		result := table.FindOne(
			ctx,
			actions.FindOne().
				Where(
					expr.Equal("ID", 1),
				),
			options.FindOne().SetDebug(true),
		)
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
		var o struct {
			Dist1 float64
			Dist2 float64
			Text  string
		}
		err = table.FindOne(
			ctx,
			actions.FindOne().
				Select(
					expr.As(expr.ST_Distance(expr.Column("Point"), origin), "dist"),
					expr.ST_Distance(
						expr.ST_GeomFromText(p1, 4326),
						expr.ST_GeomFromText(p2, 4326),
					),
					expr.ST_AsText(expr.Column("Point")),
				).
				Where(
					expr.Equal("ID", 1),
					expr.ST_Equals(origin, origin),
					// expr.ST_Within(expr.Column("Point"), orb.Point{0, 0}),
				).
				OrderBy(
					expr.Desc("dist"),
				),
			options.FindOne().SetDebug(true),
		).Scan(&o.Dist1, &o.Dist2, &o.Text)
		require.NoError(t, err)
		require.Equal(t, float64(19.6468827043885), o.Dist1)
		require.True(t, o.Dist2 > 0)
		require.Equal(t, "POINT(1 5)", o.Text)
	}
}
