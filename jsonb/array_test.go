package jsonb

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadArray(t *testing.T) {
	t.Run("with invalid value", func(t *testing.T) {
		r := NewReader([]byte(`["sdasdasd"`))
		err := r.ReadArray(func(r *Reader) error {
			return nil
		})
		require.Error(t, err)
	})

	t.Run("with null", func(t *testing.T) {
		r := NewReader([]byte(`null`))
		err := r.ReadArray(func(r *Reader) error {
			return nil
		})
		require.NoError(t, err)
	})

	t.Run("with empty array", func(t *testing.T) {
		r := NewReader([]byte(`[]`))
		err := r.ReadArray(func(r *Reader) error {
			return nil
		})
		require.NoError(t, err)
	})

	t.Run("with valid values", func(t *testing.T) {
		items := []string{}
		r := NewReader([]byte(`["[]int", "[]string" , "abc xyz"]`))
		err := r.ReadArray(func(r *Reader) error {
			items = append(items, string(r.Bytes()))
			return nil
		})
		require.NoError(t, err)
		require.ElementsMatch(t, []string{`"[]int"`, `"[]string"`, `"abc xyz"`}, items)
	})
}
