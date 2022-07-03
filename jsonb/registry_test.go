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

	{
		kind := reflect.String
		require.Panics(t, func() {
			rg.SetKindCoder(kind, textMarshalerEncoder(), nil)
		})
		require.Panics(t, func() {
			rg.SetKindCoder(kind, nil, textUnmarshalerDecoder())
		})
	}

	{
		rt := reflect.TypeOf("")
		require.Panics(t, func() {
			rg.SetTypeCoder(rt, textMarshalerEncoder(), nil)
		})
		require.Panics(t, func() {
			rg.SetTypeCoder(rt, nil, textUnmarshalerDecoder())
		})
	}

	{
		ch := make(chan any)
		v := reflect.ValueOf(ch)
		_, err = rg.LookupDecoder(v.Type())
		require.Error(t, err)
		_, err = rg.LookupEncoder(v)
		require.Error(t, err)
	}
}
