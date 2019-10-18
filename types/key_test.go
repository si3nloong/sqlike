package types

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

func TestKey(t *testing.T) {
	var (
		k   *Key
		a   bsontype.Type
		b   []byte
		err error
	)

	t.Run("ParseKey", func(it *testing.T) {
		str := `Parent,1288888/Name,'sianloong'`
		k, err = ParseKey(str)
		require.NoError(t, err)
		require.NotNil(t, k)
		require.Equal(t, NameKey("Name", "sianloong", IDKey("Parent", 1288888, nil)), k)
	})

	t.Run("MarshalBSONValue & UnmarshalBSONValue", func(it *testing.T) {
		pk := IDKey("Parent", 1288888, nil)
		require.Equal(it, "1288888", pk.ID())

		a, b, err = pk.MarshalBSONValue()
		require.NoError(it, err)
		require.Equal(it, bsontype.String, a)

		err = pk.UnmarshalBSONValue(a, b)
		require.NoError(it, err)
	})

	t.Run("MarshalJSONB & UnmarshalJSONB", func(it *testing.T) {
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

		k = new(Key)
		err = k.UnmarshalJSONB([]byte("null"))
		require.NoError(it, err)
		require.Equal(it, new(Key), k)
	})
}
