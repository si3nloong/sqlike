package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSet(t *testing.T) {

	set := Set{"a", "b", "c", "d"}

	{
		v, err := set.Value()
		require.NoError(t, err)
		require.Equal(t, "a,b,c,d", v)
	}

	{
		var set2 Set
		err := set2.Scan("a,b,c,d")
		require.NoError(t, err)
		require.Equal(t, set, set2)
	}
}
