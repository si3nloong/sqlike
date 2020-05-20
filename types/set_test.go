package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSet(t *testing.T) {

	set := Set{"a", "b", "c", "d"}

	t.Run("Set with nil value", func(it *testing.T) {
		var set Set
		v, err := set.Value()
		require.NoError(it, err)
		require.Nil(it, v)
	})

	t.Run("Set with value", func(it *testing.T) {
		v, err := set.Value()
		require.NoError(it, err)
		require.Equal(it, "a,b,c,d", v)
	})

	t.Run("Scan Set with sql.Scanner", func(it *testing.T) {
		var set2 Set
		err := set2.Scan("a,b,c,d")
		require.NoError(it, err)
		require.Equal(it, set, set2)

		set2 = Set{}
		err = set2.Scan([]byte("a,b,c,d"))
		require.NoError(t, err)
		require.Equal(t, set, set2)
	})

}
