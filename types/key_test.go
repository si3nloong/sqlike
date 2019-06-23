package types

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestKey(t *testing.T) {
	t.Run("MarshalJSONB", func(it *testing.T) {
		var (
			k   *Key
			b   []byte
			err error
		)

		pk := IDKey("Parent", 1288888, nil)
		require.Equal(it, "1288888", pk.ID())

		b, err = pk.MarshalJSONB()
		require.NoError(it, err)
		require.Equal(it, `"Parent,1288888"`, string(b))

		rk := NameKey("Name", "sianloong", pk)
		require.NoError(it, err)
		b, err = rk.MarshalJSONB()
		require.Equal(it, `"Parent,1288888/Name,'sianloong'"`, string(b))

		keyvalue := `Parent,1560407411636169424/Name,'sianloong'`
		b = []byte(strconv.Quote(keyvalue))
		err = rk.UnmarshalJSONB(b)
		require.NoError(it, err)
		require.Equal(it, keyvalue, rk.String())

		keyvalue = `Parent,1560407411636169424`
		b = []byte(strconv.Quote(keyvalue))
		k = new(Key)
		err = k.UnmarshalJSONB(b)
		require.NoError(it, err)
		require.Equal(it, keyvalue, k.String())

		keyvalue = `Parent,'a'`
		b = []byte(strconv.Quote(keyvalue))
		k = new(Key)
		err = k.UnmarshalJSONB(b)
		require.NoError(it, err)
		require.Equal(it, keyvalue, k.String())
	})
}
