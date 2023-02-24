package expr

import (
	"testing"

	"github.com/si3nloong/sqlike/v2/internal/primitive"
	"github.com/stretchr/testify/require"
)

func TestAsc(t *testing.T) {
	require.Equal(t, primitive.Sort{}, Asc("c"))
}

func TestDesc(t *testing.T) {

}
