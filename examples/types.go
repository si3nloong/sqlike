package examples

import (
	"context"
	"encoding/json"
	"sort"
	"time"

	"cloud.google.com/go/civil"
	"github.com/brianvoe/gofakeit"
	"github.com/si3nloong/sqlike/v2/db"
	"github.com/si3nloong/sqlike/v2/sql"
	"github.com/si3nloong/sqlike/v2/types"
	"golang.org/x/text/currency"
	"golang.org/x/text/language"

	"github.com/google/uuid"
)

type indexStruct struct {
	Unique string `sqlike:",unique_index"`
	ID     string `sqlike:""`
}

// Model :
type Model struct {
	CreatedAt time.Time
	UpdatedAt time.Time
}

type normalStruct struct {
	ID            uuid.UUID `sqlike:"$Key,comment=Primary key"`
	Key           *types.Key
	PtrUUID       *uuid.UUID
	VirtualColumn string `sqlike:",generated_column"`
	Date          civil.Date
	SID           string
	Emoji         string `sqlike:""`
	FullText      string
	LongStr       string  `sqlike:",longtext"`
	CustomStrType LongStr `sqlike:",size=300"`
	EmptyByte     []byte
	Byte          []byte
	Bool          bool
	priv          int
	Skip          any `sqlike:"-"`
	Int           int `sqlike:",default=100"`
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
		Key           *types.Key
		VirtualStr    string `sqlike:",virtual_column=VirtualColumn"`
		StoredStr     string `sqlike:",stored_column"`
		NestedBool    bool
		NestedNullInt *int
	}
	JSONRaw    json.RawMessage
	Map        map[string]int
	DateTime   time.Time `sqlike:",size=0"`
	Timestamp  time.Time
	Location   *time.Location
	Language   language.Tag
	Languages  []language.Tag
	Currency   currency.Unit
	Currencies []currency.Unit
	Enum       Enum      `sqlike:",enum=SUCCESS|FAILED|UNKNOWN"`
	Set        types.Set `sqlike:",set=A|B|C"`
	Model
}

type relativeNormalStruct struct {
	ID             uuid.UUID `sqlike:",primary_key"`
	NormalStructID string    `sqlike:",foreign_key=NormalStruct:ID"`
}

type simpleStruct struct {
	ID    int64 `sqlike:",auto_increment"`
	Email string
	Name  string
	Age   uint16
}

type jsonStruct struct {
	ID     int64 `sqlike:"$Key,auto_increment"`
	Text   string
	Raw    json.RawMessage
	StrArr []string
	IntArr []int
	Map    map[string]int
	Struct struct {
		StringSlice sort.StringSlice
		IntSlice    sort.IntSlice
	}
	NullableFloat *float64
}

type uuidStruct struct {
	ID string `sql:"id,uuid"`
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
	Line2 string `sqlike:",virtual_column"` // this will not work if it's embedded struct
	City  string `sqlike:",virtual_column"` // this will not work if it's embedded struct
	State string `sqlike:",virtual_column"` // this will not work if it's embedded struct
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

type CustomValue struct {
}

var (
	_ db.ColumnDataTyper = (*CustomValue)(nil)
)

// ColumnDataType :
func (c CustomValue) ColumnDataType(ctx context.Context) *sql.Column {
	f := sql.GetField(ctx)
	return &sql.Column{
		Name:     f.Name(),
		DataType: "INT",
		Type:     "INT",
	}
}

type ptrStruct struct {
	ID int64 `sqlike:"$Key,auto_increment"`
	// CustomValue   CustomValue
	NullUUID      *uuid.UUID
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
	NullLocation  *time.Location
	NullKey       *types.Key
	NullDate      *civil.Date
	NullTime      *civil.Time
	NullEnum      *Enum `sqlike:",enum=SUCCESS|FAILED|UNKNOWN"`
}

type generatedStruct struct {
	ID     string  `sqlike:"NestedID,generated_column"`
	Amount float64 `sqlike:"Amount,generated_column"`
	Nested struct {
		ID     string  `sqlike:",stored_column=NestedID"`
		Amount float64 `sqlike:",virtual_column=Amount"`
	}
	CivilDate civil.Date
	model
	Model `sqlike:"Date"`
}

type overrideStruct struct {
	generatedStruct
	ID     int64  `sqlike:",comment=Int64 ID"`      // override string ID of generatedStruct
	Amount int    `sqlike:",comment=Int Amount"`    // override string Amount of generatedStruct
	Nested string `sqlike:",comment=String Nested"` // override string Nested of generatedStruct
}

type userStatus string

const (
	userStatusActive  userStatus = "ACTIVE"
	userStatusSuspend userStatus = "SUSPEND"
)

// User :
type User struct {
	ID        int64 `sqlike:",auto_increment"`
	Name      string
	Age       int
	Status    userStatus `sqlike:",enum=ACTIVE|SUSPEND"`
	CreatedAt time.Time  `sqlike:",default=CURRENT_TIMESTAMP"`
}

// Users :
type Users []User

// Len is part of sort.Interface.
func (usrs Users) Len() int {
	return len(usrs)
}

// Swap is part of sort.Interface.
func (usrs Users) Swap(i, j int) {
	usrs[i], usrs[j] = usrs[j], usrs[i]
}

type UserAddress struct {
	ID     int64 `sqlike:",auto_increment"`
	UserID int64 `sqlike:",foreign_key=User:ID"`
}

func newNormalStruct() normalStruct {
	now := time.Now()
	ns := normalStruct{}
	// ns.Key = types.IDKey("NormalStruct", id, nil)
	ns.ID = uuid.New()
	ns.priv = 100
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
	ns.Struct.VirtualStr = gofakeit.Sentence(10)
	ns.Struct.StoredStr = `hello world!`
	ns.Struct.NestedBool = true
	ns.Date = civil.DateOf(now)
	ns.DateTime = now
	ns.Location, _ = time.LoadLocation("Asia/Kuala_Lumpur")
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
	ns.CreatedAt = now
	ns.UpdatedAt = now
	return ns
}

func newPtrStruct() ptrStruct {
	now := time.Now()
	str := `hello world`
	uid := uuid.New()
	flag := true
	b := []byte(`hello world`)
	date, _ := civil.ParseDate("2019-01-02")
	jsonByte := json.RawMessage(`{"message":"hello world"}`)
	i := 124
	i32 := int32(-603883)
	i64 := int64(-3712897389712688393)
	u8 := uint8(88)
	u64 := uint64(37128973897126)
	enum := Success
	dt := civil.DateOf(now)
	t := civil.TimeOf(now)

	ps := ptrStruct{}
	ps.NullStr = &str
	ps.NullUUID = &uid
	ps.NullByte = &b
	ps.NullBool = &flag
	ps.NullInt = &i
	ps.NullInt32 = &i32
	ps.NullInt64 = &i64
	ps.NullDate = &date
	ps.NullUint8 = &u8
	ps.NullUint64 = &u64
	ps.NullJSONRaw = &jsonByte
	ps.NullDate = &dt
	ps.NullTime = &t
	ps.NullTimestamp = &now
	ps.NullEnum = &enum
	return ps
}

func newGeneratedStruct() *generatedStruct {
	utcNow := time.Now().UTC()
	gs := &generatedStruct{}
	gs.Nested.ID = uuid.New().String()
	gs.Nested.Amount = gofakeit.Float64Range(1, 10000)
	gs.CivilDate = civil.DateOf(utcNow)
	gs.CreatedAt = utcNow
	gs.UpdatedAt = utcNow
	return gs
}
