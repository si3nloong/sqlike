package jsonb

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReader(t *testing.T) {
	var (
		r *Reader
		// err error
	)

	r = NewReader([]byte("null"))

	require.True(t, r.IsNull())

	// All data type with null should use default value
	{
		str, _ := r.reset().ReadString()
		require.Equal(t, "", str)
		str2, _ := r.reset().ReadEscapeString()
		require.Equal(t, "", str2)
		flag, _ := r.reset().ReadBoolean()
		require.Equal(t, false, flag)
		num, _ := r.reset().ReadNumber()
		require.Equal(t, json.Number("0"), num)
		b, _ := r.reset().ReadBytes()
		require.Equal(t, []byte("null"), b)
		v, _ := r.reset().ReadValue()
		require.Equal(t, nil, v)
	}
	// require.Error(t, err)
	// _, err = r.ReadBoolean()
	// require.Error(t, err)
	// require.Error(t, r.ReadArray(nil))
}
