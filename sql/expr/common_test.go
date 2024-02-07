package expr

import (
	"testing"

	"github.com/si3nloong/sqlike/v2/internal/primitive"
	"github.com/stretchr/testify/require"
)

func TestPair(t *testing.T) {
	require.Equal(t, primitive.Pair{"a", "b"}, Pair("a", "b"))
	require.Equal(t, primitive.Pair{"data_base", "table"}, Pair("data_base", "table"))
}
