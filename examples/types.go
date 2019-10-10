package examples

import (
	"encoding/json"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/si3nloong/sqlike/types"
	"golang.org/x/text/currency"
	"golang.org/x/text/language"

	uuid "github.com/google/uuid"
)

type indexStruct struct {
	Unique string `sqlike:",unique_index"`
	ID     string `sqlike:""`
}

type normalStruct struct {
	ID            uuid.UUID `sqlike:"$Key"`
	Key           *types.Key
	SID           string `sqlike:",charset=latin1"`
	Emoji         string `sqlike:""`
	FullText      string
	LongStr       string  `sqlike:",longtext"`
	CustomStrType LongStr `sqlike:",size=300"`
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
	GeoPoint      types.GeoPoint
	Struct        struct {
		VirtualStr string `sqlike:",virtual_column"`
		StoredStr  string `sqlike:",stored_column"`
		NestedBool bool   `sqlike:""`
		// NestedNullInt *int
	}
	JSONRaw json.RawMessage
	Map     map[string]int
	// GeoPoint  types.GeoPoint
	DateTime   time.Time `sqlike:",size=0"`
	Timestamp  time.Time
	Language   language.Tag
	Languages  []language.Tag
	Currency   currency.Unit
	Currencies []currency.Unit
	Enum       Enum `sqlike:",enum=SUCCESS|FAILED|UNKNOWN"`
}

type jsonStruct struct {
	ID     int64 `sqlike:"$Key,auto_increment"`
	Text   string
	Raw    json.RawMessage
	StrArr []string
	IntArr []int
	Map    map[string]int
	Struct struct {
	}
	NullableFloat *float64
	GeoPoint      types.GeoPoint
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
	Line2 string `sqlike:",virtual_column"`
	City  string `sqlike:",virtual_column"`
	State string `sqlike:",virtual_column"`
	// Country `sqlike:",inline"`
	Country Country
}

// Enum :
type Enum string

// enum :
const (
	Success Enum = "SUCCESS"
	Failed  Enum = "FAILED"
	Unknown Enum = "UNKNOWN"
)

type model struct {
	No int64
	ID string `sqlike:"id"`
	Address
}

type ptrStruct struct {
	ID            int64   `sqlike:"$Key,auto_increment"`
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
	NullKey       *types.Key
	NullDate      *types.Date
	NullEnum      *Enum `sqlike:",enum=SUCCESS|FAILED|UNKNOWN"`
}

type generatedStruct struct {
	ID     string  `sqlike:"NestedID,generated_column"`
	Amount float64 `sqlike:"Amount,generated_column"`
	Nested struct {
		ID     string  `sqlike:",stored_column=NestedID"`
		Amount float64 `sqlike:",virtual_column=Amount"`
	}
}

type mongoStruct struct {
	Key  *types.Key
	Name string
}

func newNormalStruct() normalStruct {
	now := time.Now()
	ns := normalStruct{}
	// ns.Key = types.IDKey("NormalStruct", id, nil)
	ns.ID = uuid.New()
	ns.Emoji = `😀 😁 😂 🤣 😃 😄 😅 😆 😉 😊`
	ns.Byte = []byte(`-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCklQio4TeIZo63S0FvNonY2/nA
ZUvrnDRPIzEKK4A7Hu4UjxNhebxuEA/PqSJgxOIHVPnASrSwj+IlPokcdrR6Ekyn
0cvjjwjGRyAGawVhf7TWHjkxTK6pIIqRiBK4h+E/fPwpvJTieFCSmIWovR8Wz6Jy
eCnpmNrTzG6ZJlJcvQIDAQAB
-----END PUBLIC KEY-----`)
	ns.CustomStrType = LongStr(gofakeit.RandString([]string{
		`覚えられなくて使うたびにググってしまうので、以後楽をするためにスニペットを記す。`,
		`このパッケージができた背景は`,
		`この記事ではerrorsパッケージの仕様を紹介します。`,
		`errors.Newで作成したエラーは、%+v のときにファイル名やメソッド名を表示します。`,
	}))
	ns.LongStr = gofakeit.Sentence(50)
	ns.Key = types.NewNameKey("Name", types.NewIDKey("ID", nil))
	ns.Bool = true
	ns.FullText = "Hal%o%()#$\\%^&_"
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
	ns.GeoPoint = [2]float64{0.11, 0.12312}
	ns.Struct.VirtualStr = gofakeit.Sentence(10)
	ns.Struct.StoredStr = `hello world!`
	ns.Struct.NestedBool = true
	ns.DateTime = now
	ns.Timestamp = now
	ns.Language = language.English
	ns.Currencies = []currency.Unit{
		currency.AUD,
		currency.EUR,
	}
	ns.Enum = Enum(gofakeit.RandString([]string{
		"SUCCESS",
		"FAILED",
		"UNKNOWN",
	}))
	return ns
}

func newPtrStruct() ptrStruct {
	now := time.Now()
	str := `hello world`
	flag := true
	b := []byte(`hello world`)
	date, _ := types.ParseDate("2019-01-02")
	jsonByte := json.RawMessage(`{"message":"hello world"}`)
	i := 124
	i32 := int32(-603883)
	i64 := int64(-3712897389712688393)
	u8 := uint8(88)
	u64 := uint64(37128973897126)
	enum := Success

	ps := ptrStruct{}
	ps.NullStr = &str
	ps.NullByte = &b
	ps.NullBool = &flag
	ps.NullInt = &i
	ps.NullInt32 = &i32
	ps.NullInt64 = &i64
	ps.NullDate = date
	ps.NullUint8 = &u8
	ps.NullUint64 = &u64
	ps.NullJSONRaw = &jsonByte
	ps.NullTimestamp = &now
	ps.NullEnum = &enum
	return ps
}

func newGeneratedStruct() *generatedStruct {
	gs := &generatedStruct{}
	gs.Nested.ID = uuid.New().String()
	gs.Nested.Amount = gofakeit.Float64Range(1, 10000)
	return gs
}
