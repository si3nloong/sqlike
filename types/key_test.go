package types

import (
	"encoding/json"
	"strconv"
	"testing"

	"github.com/si3nloong/sqlike/jsonb"
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

	t.Run("Clone", func(it *testing.T) {
		str := `Parent,1288888/Name,'sianloong'`
		k, err = ParseKey(str)
		require.NoError(t, err)
		require.NotNil(t, k)
		require.Equal(t, k, k.Clone())
	})

	t.Run("Encode & Decode", func(it *testing.T) {
		str := `Parent,1288888/Name,'sianloong'`
		k, err = ParseKey(str)
		require.NoError(it, err)
		require.Equal(it, `EgROYW1lIglzaWFubG9vbmcqDBIGUGFyZW50GLjVTg`, k.Encode())

		var pk *Key
		pk, err = DecodeKey(k.Encode())
		require.NoError(it, err)
		require.Equal(it, NameKey("Name", "sianloong", IDKey("Parent", 1288888, nil)), pk)
	})

	t.Run("MarshalerText & UnmarshalText", func(it *testing.T) {
		pk := IDKey("Parent", 1288888, nil)
		require.Equal(it, "1288888", pk.ID())

		b, err = pk.MarshalText()
		require.NoError(it, err)
		require.Equal(it, b, []byte(`Parent,1288888`))

		str := `EgROYW1lIg1zaWFubG9vbmcvQDkwKhISBlBhcmVudBjQ1deb4Mjr0xU`
		err = pk.UnmarshalText([]byte(str))
		require.NoError(it, err)
		require.Equal(it, "Parent,1560407411636169424/Name,'sianloong%2F@90'", pk.String())
		require.Equal(it, str, pk.Encode())
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

	t.Run("JSONB Marshal & Unmarshal", func(it *testing.T) {
		pk := IDKey("Parent", 1288888, nil)
		require.Equal(it, "1288888", pk.ID())

		b, err = pk.MarshalJSONB()
		require.NoError(it, err)
		require.Equal(it, `"Parent,1288888"`, string(b))

		rk := NameKey("Name", "sianloong", pk)
		require.NoError(it, err)
		b, err = rk.MarshalJSONB()
		require.Equal(it, `"Parent,1288888/Name,'sianloong'"`, string(b))

		kv := `Parent,1560407411636169424/Name,'sianloong'`
		b = []byte(strconv.Quote(kv))
		err = rk.UnmarshalJSONB(b)
		require.NoError(it, err)
		require.Equal(it, kv, rk.String())

		kv = `Parent,1560407411636169424`
		b = []byte(strconv.Quote(kv))
		k = new(Key)
		err = k.UnmarshalJSONB(b)
		require.NoError(it, err)
		require.Equal(it, kv, k.String())

		kv = `Parent,'a'`
		b = []byte(strconv.Quote(kv))
		k = new(Key)
		err = k.UnmarshalJSONB(b)
		require.NoError(it, err)
		require.Equal(it, kv, k.String())

		k = new(Key)
		err = k.UnmarshalJSONB([]byte("null"))
		require.NoError(it, err)
		require.Equal(it, new(Key), k)

		nk := NameKey("Name", "sianloong", pk)
		b, err = jsonb.Marshal(nk)
		require.NoError(it, err)
		require.Equal(it, `"Parent,1560407411636169424/Name,'sianloong'"`, string(b))

		k2 := new(Key)
		err = jsonb.Unmarshal(b, k2)
		require.NoError(it, err)
		require.Equal(it, nk, k2)

		k3 := new(Key)
		err = jsonb.Unmarshal([]byte(`null`), k3)
		require.NoError(it, err)
		require.Equal(it, &Key{}, k3)
	})

	t.Run("JSON Marshal & Unmarshal", func(it *testing.T) {
		k := NameKey("Name", "sianloong", nil)
		b := []byte(`"EgROYW1lIglzaWFubG9vbmc"`)
		binary, err := json.Marshal(k)
		require.NoError(it, err)
		require.Equal(it, b, binary)

		var k2 *Key
		err = json.Unmarshal(binary, &k2)
		require.NoError(it, err)
		require.Equal(it, k, k2)

		k3 := new(Key)
		err = json.Unmarshal([]byte(`null`), k3)
		require.NoError(it, err)
		require.Equal(it, &Key{}, k3)
	})
}
