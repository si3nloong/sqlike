package jsonb

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/si3nloong/sqlike/types"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/sjson"
)

type longStr string

type Decimal float64

func (f Decimal) MarshalJSONB() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%.2f"`, f)), nil
}

type Boolean bool

func (b Boolean) MarshalJSON() ([]byte, error) {
	if b {
		return []byte(`"Yes"`), nil
	}
	return []byte(`"No"`), nil
}

type Text string

func (txt Text) MarshalText() ([]byte, error) {
	return []byte(txt), nil
}

type normalStruct struct {
	Str               string
	UUID              uuid.UUID
	Text              Text
	DecimalStr        Decimal
	NullDecimalStr    *Decimal `sqlike:"NullableDecimal"`
	Boolean           Boolean
	LongStr           string
	CustomStrType     longStr
	EmptyByte         []byte
	Byte              []byte
	Bool              bool
	priv              int
	Skip              interface{}
	Int               int
	TinyInt           int8
	SmallInt          int16
	MediumInt         int32
	BigInt            int64
	Uint              uint
	TinyUint          uint8
	SmallUint         uint16
	MediumUint        uint32
	BigUint           uint64
	Float32           float32
	Float64           float64
	UFloat32          float32
	EmptyArray        []string
	EmptyMap          map[string]interface{}
	EmptyStruct       struct{}
	JSONRaw           json.RawMessage
	Timestamp         time.Time
	NullStr           *string
	NullCustomStrType *longStr
	NullInt           *int
	NullTinyInt       *int8
	NullSmallInt      *int16
	NullMediumInt     *int32
	NullBigInt        *int64
	NullFloat32       *float32
	NullFloat64       *float64
	MultiPtr          *****int
	NullKey           *types.Key
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
		i   normalStruct
		i32 = int(888)
		k   = types.NameKey("Name", "@#$%^&*()ashdkjashd", types.NewIDKey("ID", nil))
	)

	data := `{"Str":"","UUID":"00000000-0000-0000-0000-000000000000","Text":"","DecimalStr":"0.00","NullableDecimal":null,`
	data += `"Boolean":"No","LongStr":"","CustomStrType":"",`
	data += `"EmptyByte":null,"Byte":null,"Bool":false,"Skip":null,`
	data += `"Int":0,"TinyInt":0,"SmallInt":0,"MediumInt":0,"BigInt":0,`
	data += `"Uint":0,"TinyUint":0,"SmallUint":0,"MediumUint":0,"BigUint":0,`
	data += `"Float32":0,"Float64":0,"UFloat32":0,`
	data += `"EmptyArray":null,"EmptyMap":null,"EmptyStruct":{},"JSONRaw":null,"Timestamp":"0001-01-01T00:00:00Z",`
	data += `"NullStr":null,"NullCustomStrType":null,`
	data += `"NullInt":null,"NullTinyInt":null,"NullSmallInt":null,"NullMediumInt":null,"NullBigInt":null,`
	data += `"NullFloat32":null,"NullFloat64":null,`
	data += `"MultiPtr":null,"NullKey":null}`
	dataByte := []byte(data)

	// Marshal nil
	{
		b, err = Marshal(nsPtr)
		require.NoError(t, err)
		require.Equal(t, []byte(`null`), b)

		b, err = Marshal(nil)
		require.NoError(t, err)
		require.Equal(t, []byte(`null`), b)
	}

	// Marshal initialized struct
	{
		b, err = Marshal(nsInit)
		require.NoError(t, err)
		require.Equal(t, dataByte, b)
	}

	var (
		symbolStr     = `'ajhdjasd12380912$%^&*()_\\"asdasd123910293"""\\\\123210312930-\\`
		jsonEscapeStr = `"'ajhdjasd12380912$%^&*()_\\\\\"asdasd123910293\"\"\"\\\\\\\\123210312930-\\\\"`
		bytes         = []byte(`abcd1234`)
		// result        = []byte(`{"Str":"hello world","DecimalStr":"10.69","NullableDecimal":null,"Boolean":"Yes","LongStr":"` + jsonEscapeStr + `","CustomStrType":"","EmptyByte":"YWJjZDEyMzQ=","Byte":null,"Bool":false,"Skip":null,"Int":0,"TinyInt":0,"SmallInt":0,"MediumInt":0,"BigInt":0,"Uint":0,"TinyUint":0,"SmallUint":0,"MediumUint":0,"BigUint":0,"Float32":0,"Float64":0,"UFloat32":0,"EmptyStruct":{},"JSONRaw":null,"Timestamp":"0001-01-01T00:00:00Z","NullStr":"` + jsonEscapeStr + `","NullCustomStrType":null,"NullInt":888,"NullTinyInt":null,"NullSmallInt":null,"NullMediumInt":null,"NullBigInt":null,"NullFloat32":null,"NullFloat64":null,"MultiPtr":null,"NullKey":"` + k.String() + `"}`)
	)

	// Marshal struct with pointer value
	{
		i.Text = `"My long text.......""`
		i.Str = "hello world"
		i.LongStr = symbolStr
		i.Boolean = true
		i.DecimalStr = 10.688
		i.EmptyByte = bytes
		i.EmptyArray = make([]string, 0)
		i.EmptyMap = make(map[string]interface{})
		i.NullStr = &symbolStr
		i.NullInt = &i32
		i.NullKey = k

		dataByte, _ = sjson.SetBytes(dataByte, "Str", "hello world")
		dataByte, _ = sjson.SetBytes(dataByte, "Text", `"My long text.......""`)
		dataByte, _ = sjson.SetRawBytes(dataByte, "LongStr", []byte(jsonEscapeStr))
		dataByte, _ = sjson.SetBytes(dataByte, "Boolean", "Yes")
		dataByte, _ = sjson.SetBytes(dataByte, "DecimalStr", "10.69")
		dataByte, _ = sjson.SetBytes(dataByte, "EmptyByte", base64.StdEncoding.EncodeToString(bytes))
		dataByte, _ = sjson.SetBytes(dataByte, "EmptyArray", make([]string, 0))
		dataByte, _ = sjson.SetBytes(dataByte, "EmptyMap", make(map[string]interface{}))
		dataByte, _ = sjson.SetRawBytes(dataByte, "NullStr", []byte(jsonEscapeStr))
		dataByte, _ = sjson.SetBytes(dataByte, "NullInt", i32)
		dataByte, _ = sjson.SetBytes(dataByte, "NullKey", k.String())

		b, err = Marshal(i)
		require.Equal(t, dataByte, b)
	}
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
