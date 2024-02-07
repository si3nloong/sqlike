package types

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"testing"

	"github.com/si3nloong/sqlike/v2/jsonb"
	"github.com/si3nloong/sqlike/v2/sql"
	"github.com/si3nloong/sqlike/v2/x/reflext"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

type field struct {
	name string
	t    reflect.Type
	null bool
}

func (f field) Name() string {
	return f.name
}

func (f field) Type() reflect.Type {
	return f.t
}

func (field) Index() []int {
	return nil
}

func (f field) Tag() reflext.FieldTag {
	return reflext.StructTag{}
}

func (field) Parent() reflext.FieldInfo {
	return nil
}

func (field) ParentByTraversal(cb func(reflext.FieldInfo) bool) reflext.FieldInfo {
	return nil
}

func (field) Children() []reflext.FieldInfo {
	return nil
}

func (f field) IsNullable() bool {
	return f.null
}

func (field) IsEmbedded() bool {
	return false
}

var _ reflext.FieldInfo = (*field)(nil)

func TestKey(t *testing.T) {
	var (
		k   *Key
		a   bsontype.Type
		b   []byte
		err error
	)

	t.Run("DataType", func(it *testing.T) {
		k := new(Key)

		col := k.ColumnDataType(sql.Context("", "").
			SetField(field{
				name: "Key",
				t:    reflect.TypeOf(k),
			}))

		require.Equal(it, "Key", col.Name)
		require.Equal(it, "VARCHAR", col.DataType)
		require.Equal(it, "VARCHAR(512)", col.Type)
		require.Equal(it, "latin1", *col.Charset)
		require.Equal(it, "latin1_bin", *col.Collation)
		require.True(it, col.Nullable)
	})

	t.Run("String and GoString", func(it *testing.T) {
		k := NameKey("Child", "!@#$%^&*()'dajhdkhas", NameKey("Parent", "{}|askjdhkashy8q,L:\"><<", nil))

		encodedKey := `Parent,'%7B%7D%7Caskjdhkashy8q%2CL:%22%3E%3C%3C'/Child,'%21@%23$%25%5E&%2A%28%29%27dajhdkhas'`
		rawStrKey := `Parent,'{}|askjdhkashy8q,L:"><<'/Child,'!@#$%^&*()'dajhdkhas'`
		require.Equal(it, encodedKey, k.String())
		require.Equal(it, rawStrKey, k.GoString())
		require.Equal(it, rawStrKey, fmt.Sprintf("%#v", k))
	})

	t.Run("ID", func(it *testing.T) {
		nk := NameKey("Name", "name-value", nil)
		require.Equal(it, "name-value", nk.ID())

		idk := IDKey("Name", 217371238213213, nil)
		require.Equal(it, "217371238213213", idk.ID())

		var nilkey *Key
		require.Equal(it, "", nilkey.ID())
	})

	t.Run("Empty Key", func(it *testing.T) {
		k := new(Key)
		require.True(it, k.Incomplete())
		require.True(it, k.IsZero())
		require.Equal(it, k, k.Root())
		require.Nil(it, k.Parent)
	})

	t.Run("ParseKey", func(it *testing.T) {
		str := `Parent,1288888/Name,'sianloong%40gmail.com'`
		k, err = ParseKey(str)
		require.NoError(it, err)
		require.NotNil(it, k)
		nk := NameKey("Name", "sianloong@gmail.com", IDKey("Parent", 1288888, nil))
		require.Equal(it, nk, k)
		require.True(it, nk.Equal(k))
	})

	t.Run("Clone", func(it *testing.T) {
		str := `Parent,1288888/Name,'sianloong'`
		k, err = ParseKey(str)
		require.NoError(it, err)
		require.NotNil(it, k)
		require.Equal(it, k, k.Clone())
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

	t.Run("Encode & Decode Unicode", func(it *testing.T) {
		str := `Parent,1288888/Name,'%F0%9F%A4%94%E3%83%A4%E3%83%9E%E3%83%88'`
		k, err = ParseKey(str)
		require.NoError(it, err)
		require.Equal(it, `EgROYW1lIg3wn6SU44Ok44Oe44OIKgwSBlBhcmVudBi41U4`, k.Encode())

		str = `Parent,1288888/Name,'🤔ヤマト'`
		k, err = ParseKey(str)
		require.NoError(it, err)
		require.Equal(it, `EgROYW1lIg3wn6SU44Ok44Oe44OIKgwSBlBhcmVudBi41U4`, k.Encode())

		var pk *Key
		pk, err = DecodeKey(`EgROYW1lIg3wn6SU44Ok44Oe44OIKgwSBlBhcmVudBi41U4`)
		require.NoError(it, err)
		require.Equal(it, NameKey("Name", "🤔ヤマト", IDKey("Parent", 1288888, nil)), pk)
	})

	t.Run("MarshalerText & UnmarshalText", func(t *testing.T) {
		pk := IDKey("Parent", 1288888, nil)
		require.Equal(t, "1288888", pk.ID())

		b, err = pk.MarshalText()
		require.NoError(t, err)
		require.Equal(t, b, []byte(`Parent,1288888`))

		pk = new(Key)
		err = pk.UnmarshalText([]byte(``))
		require.Error(t, err)

		pk = new(Key)
		err = pk.UnmarshalText([]byte(`null`))
		require.NoError(t, err)
		require.True(t, pk.Incomplete())

		pk = new(Key)
		err = pk.UnmarshalText([]byte(`nullunknown`))
		require.Error(t, err)

		str := `EgROYW1lIg1zaWFubG9vbmcvQDkwKhISBlBhcmVudBjQ1deb4Mjr0xU`
		err = pk.UnmarshalText([]byte(str))
		require.NoError(t, err)
		require.False(t, pk.Incomplete())
		require.Equal(t, "Parent,1560407411636169424/Name,'sianloong%2F@90'", pk.String())
		require.Equal(t, str, pk.Encode())
	})

	t.Run("MarshalBSONValue & UnmarshalBSONValue", func(t *testing.T) {
		pk := IDKey("Parent", 1288888, nil)
		require.Equal(t, "1288888", pk.ID())

		a, b, err = pk.MarshalBSONValue()
		require.NoError(t, err)
		require.Equal(t, bsontype.String, a)

		k := new(Key)
		pt, _, err := k.MarshalBSONValue()
		require.NoError(t, err)
		require.Equal(t, bsontype.Null, pt)

		err = pk.UnmarshalBSONValue(a, b)
		require.NoError(t, err)

		err = pk.UnmarshalBSONValue(bsontype.Null, nil)
		require.NoError(t, err)

		var nilKey *Key
		err = nilKey.UnmarshalBSONValue(a, b)
		require.Error(t, err)
	})

	t.Run("JSONB Marshal & Unmarshal", func(t *testing.T) {
		var nilKey *Key
		b, err = jsonb.Marshal(nilKey)
		require.NoError(t, err)
		require.Equal(t, []byte(`null`), b)

		var o struct {
			Key *Key
		}

		ik := new(Key)
		b, err = json.Marshal(ik)
		require.NoError(t, err)
		require.Equal(t, []byte(`null`), b)

		ik2 := new(Key)
		b, err = jsonb.Marshal(ik2)
		require.NoError(t, err)
		require.Equal(t, []byte(`null`), b)

		b, err = jsonb.Marshal(o)
		require.NoError(t, err)
		require.Equal(t, []byte(`{"Key":null}`), b)

		pk := IDKey("Parent", 1288888, nil)
		require.Equal(t, "1288888", pk.ID())

		b, err = pk.MarshalJSONB()
		require.NoError(t, err)
		require.Equal(t, `"Parent,1288888"`, string(b))

		rk := NameKey("Name", "sianloong", pk)
		require.NoError(t, err)
		b, err = rk.MarshalJSONB()
		require.Equal(t, `"Parent,1288888/Name,'sianloong'"`, string(b))

		kv := `Parent,1560407411636169424/Name,'sianloong'`
		b = []byte(strconv.Quote(kv))

		err = rk.UnmarshalJSONB([]byte(`""`))
		require.Error(t, err)

		err = rk.UnmarshalJSONB(b)
		require.NoError(t, err)
		require.Equal(t, kv, rk.String())

		kv = `Parent,1560407411636169424`
		b = []byte(strconv.Quote(kv))
		k = new(Key)
		err = k.UnmarshalJSONB(b)
		require.NoError(t, err)
		require.Equal(t, kv, k.String())

		kv = `Parent,'a'`
		b = []byte(strconv.Quote(kv))
		k = new(Key)
		err = k.UnmarshalJSONB(b)
		require.NoError(t, err)
		require.Equal(t, kv, k.String())

		k = new(Key)
		err = k.UnmarshalJSONB([]byte("null"))
		require.NoError(t, err)
		require.Equal(t, new(Key), k)

		nk := NameKey("Name", "sianloong", pk)
		b, err = jsonb.Marshal(nk)
		require.NoError(t, err)
		require.Equal(t, `"Parent,1560407411636169424/Name,'sianloong'"`, string(b))

		k2 := new(Key)
		err = jsonb.Unmarshal(b, k2)
		require.NoError(t, err)
		require.Equal(t, nk, k2)

		k3 := new(Key)
		err = jsonb.Unmarshal([]byte(`null`), k3)
		require.NoError(t, err)
		require.Equal(t, &Key{}, k3)
	})

	t.Run("JSON Marshal & Unmarshal", func(t *testing.T) {
		k := NameKey("Name", "sianloong", nil)
		b := []byte(`"EgROYW1lIglzaWFubG9vbmc"`)
		binary, err := json.Marshal(k)
		require.NoError(t, err)
		require.Equal(t, b, binary)

		var k2 *Key
		err = json.Unmarshal([]byte(``), &k2)
		require.Error(t, err)

		err = json.Unmarshal(binary, &k2)
		require.NoError(t, err)
		require.Equal(t, k, k2)

		k3 := new(Key)
		err = json.Unmarshal([]byte(`null`), k3)
		require.NoError(t, err)
		require.Equal(t, &Key{}, k3)

		k = new(Key)
		err = json.Unmarshal([]byte(`unknown`), k)
		require.Error(t, err)

		k = new(Key)
		err = json.Unmarshal([]byte(""), k)
		require.Error(t, err)
	})

	t.Run("sql.Scanner", func(t *testing.T) {
		var k = new(Key)
		var keystr = `Parent,'hello-world%21'/Child,'hRTYUI%2CO88191'`
		err := k.Scan(keystr)
		require.NoError(t, err)
		require.Equal(t, keystr, k.String())

		err = k.Scan([]byte(keystr))
		require.NoError(t, err)
		require.Equal(t, keystr, k.String())

		err = k.Scan(`unknown`)
		require.Error(t, err)

		err = k.Scan([]byte(`wrong key`))
		require.Error(t, err)
	})

	t.Run("driver.Valuer", func(it *testing.T) {
		k := NameKey("Parent", "hello-world!", nil)
		v, err := k.Value()
		require.NoError(it, err)
		require.Equal(it, `Parent,'hello-world%21'`, v)

		nk := NameKey("Child", "hRTYUI,O88191", k)
		v, err = nk.Value()
		require.NoError(it, err)
		require.Equal(it, `Parent,'hello-world%21'/Child,'hRTYUI%2CO88191'`, v)

		idk := IDKey("Parent", 187239123213, nil)
		v, err = idk.Value()
		require.NoError(it, err)
		require.Equal(it, `Parent,187239123213`, v)

		idck := IDKey("Child", 17288, idk)
		v, err = idck.Value()
		require.NoError(it, err)
		require.Equal(it, `Parent,187239123213/Child,17288`, v)

		mk := NameKey("Mix", "Name-value", idk)
		v, err = mk.Value()
		require.NoError(it, err)
		require.Equal(it, `Parent,187239123213/Mix,'Name-value'`, v)
	})

	t.Run("MarshalGQL and UnmarshalGQL", func(it *testing.T) {
		k := NameKey("Parent", "hello-world!", nil)
		w := new(bytes.Buffer)
		k.MarshalGQL(w)
		require.NoError(it, err)
		require.Equal(it, `"EgZQYXJlbnQiDGhlbGxvLXdvcmxkIQ"`, w.String())

		nk := NameKey("Child", "hRTYUI,O88191", k)
		encoded := `"EgVDaGlsZCINaFJUWVVJLE84ODE5MSoWEgZQYXJlbnQiDGhlbGxvLXdvcmxkIQ"`
		w.Reset()
		nk.MarshalGQL(w)
		require.NoError(it, err)
		require.Equal(it, encoded, w.String())

		w.Reset()
		pk := new(Key)
		pk.MarshalGQL(w)
		require.NoError(it, err)
		require.Equal(it, `null`, w.String())

		pk = new(Key)
		err = pk.UnmarshalGQL(k)
		require.NoError(t, err)
		require.Equal(t, `Parent,'hello-world%21'`, pk.String())

		err = pk.UnmarshalGQL(encoded)
		require.NoError(t, err)
		require.Equal(t, `Parent,'hello-world%21'/Child,'hRTYUI%2CO88191'`, pk.String())

		err = pk.UnmarshalGQL([]byte(encoded))
		require.NoError(t, err)
		require.Equal(t, `Parent,'hello-world%21'/Child,'hRTYUI%2CO88191'`, pk.String())
	})

	t.Run("Check Panic", func(it *testing.T) {
		var (
			k       Key
			nilKey  *Key
			nullKey *Key
		)

		require.True(it, k.Incomplete())
		require.True(it, nilKey.Incomplete())
		require.False(it, k.Equal(nilKey))
		require.True(it, nilKey.Equal(nullKey))

		v, err := k.Value()
		require.NoError(t, err)
		require.Nil(t, v)

		require.Panics(it, func() {
			nilKey.String()
		})
		require.Panics(it, func() {
			nilKey.GoString()
		})
		require.Panics(it, func() {
			nilKey.MarshalText()
		})
		require.Panics(it, func() {
			nilKey.MarshalJSON()
		})
		require.Panics(it, func() {
			nilKey.MarshalJSONB()
		})
		require.Panics(it, func() {
			nilKey.MarshalBSONValue()
		})
		require.Panics(it, func() {
			nilKey.GobEncode()
		})
		require.Panics(it, func() {
			nilKey.Encode()
		})
	})

	t.Run("keyToGobKey", func(t *testing.T) {
		var k *Key
		gk := keyToGobKey(k)
		require.Nil(t, gk)
	})

	t.Run("gobKeyToKey", func(t *testing.T) {
		var gk *gobKey
		k := gobKeyToKey(gk)
		require.Nil(t, k)
	})

	nk := NewNameKey("Name", nil)
	require.NotEmpty(t, nk.NameID)
	require.Equal(t, "Name", nk.Kind)
	require.Empty(t, nk.IntID)
	require.Nil(t, nk.Parent)

	idk := NewIDKey("ID", nil)
	require.Empty(t, idk.NameID)
	require.Equal(t, "ID", idk.Kind)
	require.NotEmpty(t, idk.IntID)
	require.Nil(t, idk.Parent)
}
