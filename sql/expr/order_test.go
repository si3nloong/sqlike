package expr

import (
	"testing"

	"github.com/si3nloong/sqlike/v2/internal/primitive"
	"github.com/stretchr/testify/require"
)

func TestAsc(t *testing.T) {
	require.Equal(t, primitive.Sort{Field: wrapColumn("c"), Order: primitive.Ascending}, Asc("c"))
}

func TestDesc(t *testing.T) {
	require.Equal(t, primitive.Sort{Field: wrapColumn("c"), Order: primitive.Descending}, Desc("c"))
}
