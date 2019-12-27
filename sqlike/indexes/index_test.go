package indexes

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIndex(t *testing.T) {
	idx := Index{Columns: Columns("a", "b2", "a_c3", "d.h8")}
	require.Equal(t, `587bc84ba16ffe5618f4864bcea6c9a6`, idx.GetName())
	require.Equal(t, `587bc84ba16ffe5618f4864bcea6c9a6`, idx.HashName())
}
