package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStringSlice(t *testing.T) {
	slice := StringSlice{"z", "a", "e", "abc", "b", "c", "dd"}

	require.Equal(t, -1, slice.IndexOf("cc"))
	require.Equal(t, 1, slice.IndexOf("a"))
	require.Equal(t, 4, slice.IndexOf("b"))
	require.Equal(t, 5, slice.IndexOf("c"))

	slice.Sort()
	require.ElementsMatch(t, StringSlice{"a", "abc", "b", "c", "dd", "e", "z"}, slice)

	slice.Splice(1)
	require.ElementsMatch(t, StringSlice{"a", "b", "c", "dd", "e", "z"}, slice)

	slice.Splice(3)
	require.ElementsMatch(t, StringSlice{"a", "b", "c", "e", "z"}, slice)
}
