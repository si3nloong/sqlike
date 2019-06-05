package jsonb

import (
	"encoding/json"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type longStr string

type normalStruct struct {
	Str           string
	LongStr       string
	CustomStrType longStr
	EmptyByte     []byte
	Byte          []byte
	Bool          bool
	priv          int
	Skip          interface{}
	Int           int
	TinyInt       int8
	SmallInt      int16
	MediumInt     int32
	BigInt        int64
	Uint          uint
	TinyUint      uint8
	SmallUint     uint16
	MediumUint    uint32
	BigUint       uint64
	Float32       float32
	Float64       float64
	UFloat32      float32
	EmptyStruct   struct{}
	JSONRaw       json.RawMessage
	Timestamp     time.Time
}

var (
	nsPtr  *normalStruct
	ns     normalStruct
	nsInit = new(normalStruct)
)

func TestMarshal(t *testing.T) {
	var (
		b   []byte
		err error
	)

	b, err = json.Marshal(nsPtr)
	log.Println(string(b))
	require.NoError(t, err)

	b, err = json.Marshal(nsInit)
	log.Println(string(b))
	require.NoError(t, err)

	b, err = json.Marshal(nsPtr)
	log.Println(string(b))
	require.NoError(t, err)

	b, err = Marshal(nsPtr)
	log.Println(string(b))
	require.NoError(t, err)

	b, err = Marshal(nsInit)
	log.Println(string(b))
	require.NoError(t, err)

	b, err = Marshal(nsPtr)
	log.Println(string(b))
	require.NoError(t, err)

	// output := `{"Str":"","LongStr":"","CustomStrType":"",`
	// output += `"EmptyByte":null,"Byte":null,"Bool":false,`
	// output += `"Int":0,"TinyInt":0,"SmallInt":0,"MediumInt":0,`
	// output += `"BigInt":0,"Uint":0,"TinyUint":0,"SmallUint":0,`
	// output += `"MediumUint":0,"BigUint":0,"Float32":0,"Float64":0,`
	// output += `"UFloat32":0,"EmptyStruct":{},"JSONRaw":null,`
	// output += `"Timestamp":"0001-01-01T00:00:00Z"}`
	// ins := new(normalStruct)
	// b, _ = Marshal(ins)
	// assert.Equal(t, b, []byte(output), "it should match the expected result")
}

func BenchmarkJSONMarshal(b *testing.B) {
	b.Run("Pointer Struct w/o initialize", func(t *testing.B) {
		for n := 0; n < t.N; n++ {
			json.Marshal(nsPtr)
		}
	})
	b.Run("Pointer Struct w initialize", func(t *testing.B) {
		for n := 0; n < t.N; n++ {
			json.Marshal(nsInit)
		}
	})
	b.Run("Struct w initialize", func(t *testing.B) {
		for n := 0; n < t.N; n++ {
			json.Marshal(nsPtr)
		}
	})
}

func BenchmarkJSONBMarshal(b *testing.B) {
	b.Run("Pointer Struct w/o initialize", func(t *testing.B) {
		for n := 0; n < t.N; n++ {
			Marshal(nsPtr)
		}
	})
	b.Run("Pointer Struct w initialize", func(t *testing.B) {
		for n := 0; n < t.N; n++ {
			Marshal(nsInit)
		}
	})
	b.Run("Struct w initialize", func(t *testing.B) {
		for n := 0; n < t.N; n++ {
			Marshal(nsPtr)
		}
	})
}
