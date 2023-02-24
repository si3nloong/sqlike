package jsonb

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRegistry(t *testing.T) {
	var (
		rg  = buildDefaultRegistry()
		err error
	)

	t.Run("SetKindCoder", func(t *testing.T) {
		kind := reflect.String
		require.Panics(t, func() {
			rg.SetKindCoder(kind, textMarshalerEncoder(), nil)
		})
		require.Panics(t, func() {
			rg.SetKindCoder(kind, nil, textUnmarshalerDecoder())
		})
	})

	t.Run("SetTypeCoder", func(t *testing.T) {
		rt := reflect.TypeOf("")
		require.Panics(t, func() {
			rg.SetTypeCoder(rt, textMarshalerEncoder(), nil)
		})
		require.Panics(t, func() {
			rg.SetTypeCoder(rt, nil, textUnmarshalerDecoder())
		})
	})

	t.Run("LookupDecoder", func(t *testing.T) {
		t.Run("on valid type", func(t *testing.T) {
			var str string
			v := reflect.ValueOf(str)

			dec, err := rg.LookupDecoder(v.Type())
			require.NoError(t, err)
			require.NotNil(t, dec)

			enc, err := rg.LookupEncoder(v)
			require.NoError(t, err)
			require.NotNil(t, enc)
		})

		t.Run("on missing type", func(t *testing.T) {
			ch := make(chan interface{})
			v := reflect.ValueOf(ch)
			_, err = rg.LookupDecoder(v.Type())
			require.Error(t, err)
			_, err = rg.LookupEncoder(v)
			require.Error(t, err)
		})
	})
}
