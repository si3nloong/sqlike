package examples

import (
	"encoding/json"
	"time"

	"github.com/brianvoe/gofakeit"

	uuid "github.com/satori/go.uuid"
)

type normalStruct struct {
	ID            uuid.UUID `sqlike:"$Key"`
	Emoji         string    `sqlike:""`
	LongStr       string    `sqlike:",longtext"`
	CustomStrType LongStr   `sqlike:",size:300"`
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
	UFloat32      float32 `sqlike:",unsigned"`
	EmptyStruct   struct{}
	Struct        struct {
		VirtualStr string `sqlike:",virtual_column"`
		StoredStr  string `sqlike:",stored_column"`
		NestedBool bool   `sqlike:""`
		// NestedNullInt *int
	}
	JSONRaw   json.RawMessage
	Timestamp time.Time
	Enum      Enum `sqlike:",enum:SUCCESS|FAILED|UNKNOWN"`
}

// LongStr :
type LongStr string

// Country :
type Country struct {
	Name LongStr `sqlike:""`
	Code string  `sqlike:""`
}

// Address :
type Address struct {
	Line1 string
	Line2 string `sqlike:",virtual"`
	City  string `sqlike:",virtual"`
	State string `sqlike:",virtual"`
	// Country `sqlike:",inline"`
	Country Country
}

// Enum :
type Enum string

type model struct {
	No int64
	ID string `sqlike:"id"`
	Address
}

type ptrStruct struct {
	NullStr       *string `sqlike:"nullstr"`
	NullBool      *bool
	NullByte      *[]byte
	NullInt       *int
	NullInt8      *int8
	NullInt16     *int16
	NullInt32     *int32
	NullInt64     *int64
	NullUint      *uint
	NullUint8     *uint8
	NullUint16    *uint16
	NullUint32    *uint32
	NullUint64    *uint64
	NullUFloat    *float32 `sqlike:",unsigned"`
	NullFloat32   *float32
	NullFloat64   *float64
	NullStruct    *struct{}
	NullJSONRaw   *json.RawMessage
	NullTimestamp *time.Time
	NullEnum      *Enum `sqlike:",enum:SUCCESS|FAILED|UNKNOWN"`
}

func newNormalStruct() normalStruct {
	ns := normalStruct{}
	// ns.Key = types.IDKey("NormalStruct", id, nil)
	ns.ID = uuid.NewV1()
	ns.Emoji = `😀 😁 😂 🤣 😃 😄 😅 😆 😉 😊 `
	ns.Byte = []byte(`-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCklQio4TeIZo63S0FvNonY2/nA
ZUvrnDRPIzEKK4A7Hu4UjxNhebxuEA/PqSJgxOIHVPnASrSwj+IlPokcdrR6Ekyn
0cvjjwjGRyAGawVhf7TWHjkxTK6pIIqRiBK4h+E/fPwpvJTieFCSmIWovR8Wz6Jy
eCnpmNrTzG6ZJlJcvQIDAQAB
-----END PUBLIC KEY-----`)
	ns.CustomStrType = LongStr(gofakeit.RandString([]string{
		`覚えられなくて使うたびにググってしまうので、以後楽をするためにスニペットを記す。`,
		`このパッケージができた背景は`,
		`この記事ではxerrorsパッケージの仕様を紹介します。`,
		`xerrors.Newで作成したエラーは、%+v のときにファイル名やメソッド名を表示します。`,
	}))
	ns.LongStr = gofakeit.Sentence(50)
	ns.Bool = true
	ns.Int = gofakeit.Number(100, 99999999)
	ns.TinyInt = 99
	ns.SmallInt = gofakeit.Int16()
	ns.MediumInt = gofakeit.Int32()
	ns.BigInt = gofakeit.Int64()
	ns.TinyUint = gofakeit.Uint8()
	ns.SmallUint = gofakeit.Uint16()
	ns.MediumUint = gofakeit.Uint32()
	ns.Uint = uint(gofakeit.Number(100, 99999999))
	ns.BigUint = gofakeit.Uint64()
	ns.UFloat32 = gofakeit.Float32Range(10, 10000)
	ns.Float32 = gofakeit.Float32()
	ns.Float64 = gofakeit.Float64()
	ns.JSONRaw = json.RawMessage(`{
		"message" :  "hello world",
		"code":      200,
		"error": {
			"code": "Unauthorised",
			"message": "please contact our support"
		}
	}`)
	ns.Struct.VirtualStr = gofakeit.Sentence(10)
	ns.Struct.StoredStr = `hello world!`
	ns.Struct.NestedBool = true
	ns.Timestamp = time.Now()
	ns.Enum = Enum(gofakeit.RandString([]string{
		"SUCCESS",
		"FAILED",
		"UNKNOWN",
	}))
	return ns
}
