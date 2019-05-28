package jsonb

import (
	"encoding/json"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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
	Skip          interface{} `sqlike:"-"`
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
		ns *normalStruct
		b  []byte
	)

	str := `hello world`
	b, _ = Marshal(str)
	assert.Equal(t, string(b), strconv.Quote(str), "it should be quote string")

	b, _ = Marshal(ns)
	assert.Equal(t, b, jsonNull, "it should be null")

	output := `{"Str":"","LongStr":"","CustomStrType":"",`
	output += `"EmptyByte":null,"Byte":null,"Bool":false,`
	output += `"Int":0,"TinyInt":0,"SmallInt":0,"MediumInt":0,`
	output += `"BigInt":0,"Uint":0,"TinyUint":0,"SmallUint":0,`
	output += `"MediumUint":0,"BigUint":0,"Float32":0,"Float64":0,`
	output += `"UFloat32":0,"EmptyStruct":{},"JSONRaw":null,`
	output += `"Timestamp":"0001-01-01T00:00:00Z"}`
	ins := new(normalStruct)
	b, _ = Marshal(ins)
	assert.Equal(t, b, []byte(output), "it should match the expected result")
}

func BenchmarkJSONMarshal(t *testing.B) {
	t.Run("Pointer Struct w/o initialize", func(_ *testing.B) {
		Marshal(nsPtr)
	})
	t.Run("Pointer Struct w initialize", func(_ *testing.B) {
		Marshal(nsInit)
	})
	t.Run("Struct w initialize", func(_ *testing.B) {
		Marshal(nsPtr)
	})
}

func BenchmarkJSONBMarshal(t *testing.B) {
	t.Run("Pointer Struct w/o initialize", func(_ *testing.B) {
		Marshal(nsPtr)
	})
	t.Run("Pointer Struct w initialize", func(_ *testing.B) {
		Marshal(nsInit)
	})
	t.Run("Struct w initialize", func(_ *testing.B) {
		Marshal(nsPtr)
	})
}
