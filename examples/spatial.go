package examples

import (
	"testing"

	"github.com/paulmach/orb"
	"github.com/si3nloong/sqlike/sqlike"
	"github.com/si3nloong/sqlike/sqlike/options"
	"github.com/stretchr/testify/require"
)

type Spatial struct {
	Point      orb.Point
	LineString orb.LineString
	// Polygon orb.Polygon
}

// SpatialExamples :
func SpatialExamples(t *testing.T, db *sqlike.Database) {
	var (
		sp    = Spatial{}
		table = db.Table("spatial")
		err   error
	)

	{
		err = table.DropIfExits()
		require.NoError(t, err)
	}

	{
		table.MustMigrate(Spatial{})

	}

	{
		sp.Point = orb.Point{1, 1}
		sp.LineString = []orb.Point{
			orb.Point{0, 0},
			orb.Point{1, 1},
			orb.Point{2, 2},
		}
		_, err = table.InsertOne(&sp,
			options.InsertOne().SetDebug(true),
		)
		require.NoError(t, err)
	}
}
