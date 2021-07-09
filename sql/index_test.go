package sql

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIndex(t *testing.T) {
	require.Equal(t, "BTREE", BTree.String())
	require.Equal(t, "PRIMARY", Primary.String())
	require.Equal(t, "FULLTEXT", FullText.String())
	require.Equal(t, "UNIQUE", Unique.String())
	require.Equal(t, "MULTI-VALUED", MultiValued.String())
	require.Equal(t, "SPATIAL", Spatial.String())

	idx := Index{Columns: IndexedColumns("a", "b2", "a_c3", "d.h8")}
	require.Equal(t, `587bc84ba16ffe5618f4864bcea6c9a6`, idx.GetName())
	require.Equal(t, `587bc84ba16ffe5618f4864bcea6c9a6`, idx.HashName())
}
