package mysql

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEncodeMap(t *testing.T) {
	var (
		ms  = New()
		enc = DefaultEncoders{}
	)

	t.Run("EncodeMap with nil map", func(t *testing.T) {
		var nilmap map[string]string
		query, args, err := enc.EncodeMap(ms, reflect.ValueOf(nilmap), nil)
		require.NoError(t, err)
		require.Equal(t, "?", query)
		require.ElementsMatch(t, []any{`null`}, args)
	})

	// t.Run("EncodeMap with map[string]int", func(t *testing.T) {
	// 	intmap := map[string]int{
	// 		"one":             1,
	// 		"eleven":          11,
	// 		"hundred-and-ten": 110,
	// 	}
	// 	query, args, err := enc.EncodeMap(ms, reflect.ValueOf(intmap), nil)
	// 	require.NoError(t, err)
	// 	require.Equal(t, "?", query)
	// 	require.ElementsMatch(t, []any{[]byte(`{"one":1,"eleven":11,"hundred-and-ten":110}`)}, args)
	// })

	t.Run("EncodeMap with map[string]any", func(t *testing.T) {
		intmap := make(map[string]any)
		query, args, err := enc.EncodeMap(ms, reflect.ValueOf(intmap), nil)
		require.NoError(t, err)
		require.Equal(t, "?", query)
		require.ElementsMatch(t, []any{[]byte(`{}`)}, args)
	})

	t.Run("EncodeMap with map[int]any", func(t *testing.T) {
		intmap := make(map[int]any)
		query, args, err := enc.EncodeMap(ms, reflect.ValueOf(intmap), nil)
		require.NoError(t, err)
		require.Equal(t, "?", query)
		require.ElementsMatch(t, []any{[]byte(`{}`)}, args)
	})
}
