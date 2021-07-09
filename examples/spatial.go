package examples

import (
	"context"
	"database/sql"
	"testing"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/encoding/wkb"
	"github.com/si3nloong/sqlike/v2"
	"github.com/si3nloong/sqlike/v2/actions"
	"github.com/si3nloong/sqlike/v2/options"
	sqlx "github.com/si3nloong/sqlike/v2/sql"
	"github.com/si3nloong/sqlike/v2/sql/expr"
	"github.com/stretchr/testify/require"
)

// Spatial :
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
func SpatialExamples(ctx context.Context, t *testing.T, db *sqlike.Database) {
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

	// create spatial index
	{
		table.MustMigrate(ctx, Spatial{})
		table.MustUnsafeMigrate(ctx, Spatial{})
		iv := table.Indexes()
		idx := sqlx.Index{
			Type:    sqlx.Spatial,
			Columns: sqlx.IndexedColumns("Point"),
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

	// insert spatial record
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

	// find spatial record
	{
		result := table.FindOne(
			ctx,
			actions.FindOne().
				Select(
					"ID",
					"Point",
					"PtrPoint",
					"LineString",
					"PtrLineString",
				).
				Where(
					expr.Equal("ID", 1),
				),
			options.FindOne().SetDebug(true),
		)

		var (
			c1 = new(sql.RawBytes)
			c2 = orb.Point{}
			c3 = orb.Point{}
			c4 = orb.LineString{}
			c5 *orb.LineString
		)

		cols := result.Columns()
		require.ElementsMatch(t, []string{
			"ID",
			"Point",
			"PtrPoint",
			"LineString",
			"PtrLineString",
		}, cols)
		err = result.Scan(c1, wkb.Scanner(&c2), wkb.Scanner(&c3), &c4, &c5)
		require.NoError(t, err)

		v1 := sql.RawBytes(`1`)
		nilLineString := orb.LineString(nil)
		require.Equal(t, &v1, c1)
		require.Equal(t, &v1, c1)
		require.Equal(t, orb.Point{}, c3)
		require.Equal(t, orb.LineString{{0, 0}, {1, 1}}, c4)
		require.Equal(t, &nilLineString, c5)
	}

	// find spatial record and verify the output
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

	// get distance between two point
	{
		origin := orb.Point{20, 10}
		p1 := orb.Point{1, 3}
		p2 := orb.Point{4, 18}
		var o struct {
			Dist1 float64
			Dist2 float64
			Text  string
		}
		/*
			SELECT
				ST_Distance(`Point`,ST_PointFromText("POINT(20 10)")) AS `dist`,
				ST_Distance(ST_GeomFromText("POINT(1 3)",4326),ST_GeomFromText("POINT(4 18)",4326)),ST_AsText(`Point`)
			FROM `sqlike`.`spatial`
			WHERE (
				`ID` = 1 AND
				ST_Equals(ST_PointFromText("POINT(20 10)"),ST_PointFromText("POINT(20 10)"))
			)
			ORDER BY `dist` DESC LIMIT 1;
		*/
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
