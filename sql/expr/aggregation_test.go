package expr

import (
	"testing"

	"github.com/si3nloong/sqlike/v2/internal/primitive"
	"github.com/stretchr/testify/require"
)

func TestSum(t *testing.T) {
	require.Equal(t, primitive.Aggregate{
		Field: wrapColumn("a"),
		By:    primitive.Sum,
	}, Sum("a"))
}
