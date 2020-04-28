package types

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAny(t *testing.T) {
	var (
		b   []byte
		err error
		any Any
	)

	raw := []byte(`128`)
	rawByte := []byte(`{"kind":0,"raw":128}`)

	any.raw = raw
	b, err = json.Marshal(any)
	require.NoError(t, err)

	require.Equal(t, "128", any.String())
	require.Equal(t, int64(128), any.Int64())
	require.Equal(t, uint64(128), any.Uint64())
	require.Equal(t, float64(128), any.Float64())
	require.Equal(t, rawByte, b)

	var output Any
	err = json.Unmarshal(b, &output)
	require.NoError(t, err)

	require.Equal(t, Any{kind: reflect.Invalid, raw: raw}, output)
}
