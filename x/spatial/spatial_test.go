package spatial

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSpatialType(t *testing.T) {

	require.Equal(t, "ST_GeomFromText", SpatialTypeGeomFromText.String())
	require.Equal(t, "ST_Distance", SpatialTypeDistance.String())
	require.Equal(t, "ST_Within", SpatialTypeWithin.String())
	require.Equal(t, "ST_Equals", SpatialTypeEquals.String())
	require.Equal(t, "ST_PointFromText", SpatialTypePointFromText.String())
	require.Equal(t, "ST_LineString", SpatialTypeLineString.String())
	require.Equal(t, "ST_AsText", SpatialTypeAsText.String())
	require.Equal(t, "ST_AsWKB", SpatialTypeAsWKB.String())
	require.Equal(t, "ST_AsWKT", SpatialTypeAsWKT.String())
	require.Equal(t, "ST_SRID", SpatialTypeSRID.String())
	require.Equal(t, "ST_IsValid", SpatialTypeIsValid.String())
	require.Equal(t, "ST_X", SpatialTypeX.String())
	require.Equal(t, "ST_Y", SpatialTypeY.String())
	require.Equal(t, "ST_AsGeoJSON", SpatialTypeAsGeoJSON.String())
	require.Equal(t, "ST_Area", SpatialTypeArea.String())
	require.Equal(t, "ST_Intersects", SpatialTypeIntersects.String())
	require.Equal(t, "ST_Transform", SpatialTypeTransform.String())

}
