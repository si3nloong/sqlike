package jsonb

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/si3nloong/sqlike/types"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/currency"
	"golang.org/x/text/language"
)

type User struct {
	Name  string
	Email string
	Age   int
	UUID  uuid.UUID
}

type ptrStruct struct {
	PtrStr          *string
	PtrByte         *[]byte
	PtrBool         *bool
	PtrInt          *int
	PtrInt8         *int8
	PtrInt16        *int16
	PtrInt32        *int32
	PtrInt64        *int64
	PtrUint         *uint
	PtrUint8        *uint8
	PtrUint16       *uint16
	PtrUint32       *uint32
	PtrUint64       *uint64
	PtrFloat32      *float32
	PtrFloat64      *float64
	PtrJSONRaw      *json.RawMessage
	PtrJSONNumber   *json.Number
	PtrStrSlice     *[]string
	PtrIntSlice     *[]int
	PtrUintSlice    *[]uint
	PtrF32Slice     *[]float32
	PtrF64Slice     *[]float64
	PtrLanguageTag  *language.Tag
	PtrCurrencyUnit *currency.Unit
	PtrUUID         *uuid.UUID
	PtrTime         *time.Time
	PtrStruct       *struct {
		Nested string
	}
}

type testStruct struct {
	Str        string
	UUID       *****uuid.UUID
	BigDecimal float64
	SymbolStr  string
	EscapeStr  string
	StrSlice   []string
	Users      []User
	Nested     struct {
		NilKey         *types.Key
		MultiNilPtrKey ******types.Key
		Key            *types.Key
		MultiPtr       ***string
		Security       struct {
			PrivateKey []byte
		}
		Test   string
		Test2  string
		Nested struct {
			Deep struct {
				Value float64 `sqlike:"value"`
			} `sqlike:"deep1"`
			YOLO string `sqlike:"deep2"`
		} `sqlike:"nested"`
	}
	EmptyStruct struct{}
	Integer     int
	Bool        bool
}

type customKey struct {
	value string
}

func (c *customKey) UnmarshalText(b []byte) error {
	*c = customKey{value: string(b)}
	return nil
}

func TestUnmarshal(t *testing.T) {
	var (
		strval       = `hello world !!!!!!`
		symbolstrval = `\n\t\t\t<html>\n\t\t\t\t<div>Hello World</div>\n\t\t\t</html>\n\t\t\ttesting with !@#$%^&*(_)()()_()_((*??D|}A||||\\\\))\n\t\t`
		byteval      = []byte(`LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlDWFFJQkFBS0JnUUMrYUlTemtOdXFiZmdxWW9IYW1iS0dyaEF6UnV0dWYydWFzOUJIeXllUFJUdUk5ZVdwCnJHY3lRZlhPVlh2OGJBZVMxK2tIS0MvK1ZDTk9EbGZBTFlQWVVZa053eHVvRnFMbU1SR3E1MzMwSEVLSUpySDcKSUU5aUs0QUVZL3h5WjBTUEp5ZkNnQ2ZaeGtJTmpacWFoSS8rVWxrL1BmdWwyaEQ0ZTNUZVpGTm5HUUlEQVFBQgpBb0dBSHpOYlExMWlVV3dSdFVoTkJQZ1lnKzh6NG1NbG93YW9LRUo4eDdibmRaZWZxTks2WG5KTXVyU0tSZFJHCks5ZTc2ZmtOUzBudmkxcFlLcXM0LzltMWQ4Mk9XdmtDeXZvR3pmRXdyNGJ6bVBBZjdkczVkWElhb29wbWV4WWwKbFpsSmtuMDhWNFJmOWc4RFEyNFRsb3BpZ3RrSzY5UktRSzFHaHVyV1A4UjVxeTBDUVFEZ3dxcGVKZFF5RUdaYgpQUElJN2ZsSUVDRjlQNnNVc1ovNW40cEhNNmg2N0dOdGU1cEx4bDkzOXYybVhaN09aSUZHQU1rUmNDL1ZIK3c4Cm5oaytaNE9yQWtFQTJOK01oOWRpN1YreldaNUNIWXBWTHM5Qi9xOVl3YjFCNjN0UnZUbG9QSnFqTHc1NDUzZUUKbEs0ZnJSaVhXbEhLaUpLYlBOTU1ZUVkyTVRrcEQ2dDhTd0pCQUlkU2JRVFdQZFlPcmJITkZlUnVjeUlDSkVlbQpwN2lENFUrSDBOZGhzTlNoc3BOZVVkM0JpQVZRZmhOR1ZyRHBMalFaa1BXZzJBdTNkcUpnaGM1ZXdKVUNRQUVFCkV4RnoxZGZNMGZkQ2dZYkg1aHhCQmtzZUlTbFBMS2JndmdKSDZaQVhIVnFVRThicHpXb3c0cDhaOVdPTDdJbjEKUGRyc0ZpdkNMckRPVnIzbkRMOENRUURKSENwSEVFNTc0ckpzblJNYk5rR0F5dmZheW9MeElhVUF5WXovaGxrMgpzQ0wzb3BsdDNYM0tjYzV1MkRISVFsZTdGM1M4Wmp4REZMSVRrbnJ4QS9UVgotLS0tLUVORCBSU0EgUFJJVkFURSBLRVktLS0tLQ==`)
		nullval      = []byte(`null`)
	)

	var (
		pk  []byte
		b   []byte
		err error
	)

	t.Run("Unmarshal RawMessage", func(it *testing.T) {
		var (
			b = []byte(`{
				"Raw": {
					"Title": "Header",
					"Body": {}
				},
				"Nested" : {
					"DeepNested" : {
						"RawObject" : {
							"Name" : "john",
							"Email" : "john@hotmail.com",
							"Emoji":  " 😆 😉 😊 ",
							"Age" : 28
						},
						"RawStr" : "hello world!",
						"RawNumber": 189202.882
					}
				}
			}`)
			o struct {
				Raw    json.RawMessage
				Nested struct {
					DeepNested struct {
						RawObject json.RawMessage
						RawStr    json.RawMessage
						RawNumber json.RawMessage
					}
				}
			}
		)

		err = Unmarshal(b, &o)
		require.NoError(it, err)
		require.Equal(it, json.RawMessage(`{"Title":"Header","Body":{}}`), o.Raw)
		require.Equal(it, json.RawMessage(`{"Name":"john","Email":"john@hotmail.com","Emoji":" 😆 😉 😊 ","Age":28}`), o.Nested.DeepNested.RawObject)
		require.Equal(it, json.RawMessage(`"hello world!"`), o.Nested.DeepNested.RawStr)
		require.Equal(it, json.RawMessage(`189202.882`), o.Nested.DeepNested.RawNumber)
		require.Equal(it, json.RawMessage(`"hello world!"`), o.Nested.DeepNested.RawStr)
	})

	{
		pk, err = ioutil.ReadFile("./../examples/pk.pem")
		require.NoError(t, err)
		b, err = base64.StdEncoding.DecodeString(string(byteval))
		require.NoError(t, err)
		require.Equal(t, b, pk)
	}

	t.Run("Unmarshal UUID", func(it *testing.T) {
		var uid uuid.UUID
		err = Unmarshal([]byte(`"4c03d1de-645b-40d2-9ed5-12bb537a602e"`), &uid)
		require.NoError(t, err)
		require.Equal(t, uuid.MustParse("4c03d1de-645b-40d2-9ed5-12bb537a602e"), uid)

		var ptruid *****uuid.UUID
		err = Unmarshal([]byte(`"4c03d1de-645b-40d2-9ed5-12bb537a602e"`), &ptruid)
		require.NoError(t, err)
		require.NotNil(t, ptruid)
		require.Equal(t, uuid.MustParse("4c03d1de-645b-40d2-9ed5-12bb537a602e"), *****ptruid)

		var nilUUID *uuid.UUID
		err = Unmarshal([]byte(`null`), &nilUUID)
		require.NoError(t, err)
		require.Nil(t, nilUUID)
	})

	t.Run("Unmarshal String", func(it *testing.T) {
		var (
			addrptr *string
			nilptr  *string
		)

		err = Unmarshal([]byte(`null`), &addrptr)
		require.NoError(t, err)
		require.Equal(t, nilptr, addrptr)

		var str string
		err = Unmarshal([]byte(`"`+strval+`"`), &str)
		require.Equal(it, strval, str)
		require.NoError(t, err)

		output := `
			<html>
				<div>Hello World</div>
			</html>
			testing with !@#$%^&*(_)()()_()_((*??D|}A||||\\))
		`

		var symbolstr string
		err = Unmarshal([]byte(`"`+symbolstrval+`"`), &symbolstr)
		require.NoError(t, err)
		require.Equal(it, output, symbolstr)

		err = Unmarshal(nullval, &str)
		require.NoError(t, err)
		require.Equal(it, "", str)

		var uinitstr *string
		err = Unmarshal([]byte(`null`), uinitstr)
		require.Error(t, err)

		err = Unmarshal([]byte(`null`), nil)
		require.Error(t, err)
	})

	t.Run("Unmarshal Boolean", func(it *testing.T) {
		var flag bool
		err = Unmarshal([]byte(`true`), &flag)
		require.NoError(t, err)
		require.Equal(it, true, flag)

		err = Unmarshal([]byte(`false`), &flag)
		require.NoError(t, err)
		require.Equal(it, false, flag)

		err = Unmarshal(nullval, &flag)
		require.NoError(it, err)
		require.Equal(it, false, flag)
	})

	t.Run("Unmarshal Integer", func(it *testing.T) {
		var (
			i8  int8
			i16 int16
			i32 int32
			i64 int64
			i   int
		)

		err = Unmarshal([]byte(`10`), &i8)
		require.NoError(t, err)
		require.Equal(it, int8(10), i8)

		err = Unmarshal([]byte(`-10`), &i8)
		require.NoError(t, err)
		require.Equal(it, int8(-10), i8)

		err = Unmarshal(nullval, &i8)
		require.NoError(t, err)
		require.Equal(it, int8(0), i8)

		err = Unmarshal([]byte(`128`), &i16)
		require.NoError(t, err)
		require.Equal(it, int16(128), i16)

		err = Unmarshal([]byte(`-128`), &i16)
		require.NoError(t, err)
		require.Equal(it, int16(-128), i16)

		err = Unmarshal(nullval, &i16)
		require.NoError(t, err)
		require.Equal(it, int16(0), i16)

		err = Unmarshal([]byte(`1354677198`), &i32)
		require.NoError(t, err)
		require.Equal(it, int32(1354677198), i32)

		err = Unmarshal([]byte(`-1354677198`), &i32)
		require.NoError(t, err)
		require.Equal(it, int32(-1354677198), i32)

		err = Unmarshal(nullval, &i32)
		require.NoError(t, err)
		require.Equal(it, int32(0), i32)

		err = Unmarshal([]byte(`7354673213123121983`), &i64)
		require.NoError(t, err)
		require.Equal(it, int64(7354673213123121983), i64)

		err = Unmarshal([]byte(`-7354673213123121983`), &i64)
		require.NoError(t, err)
		require.Equal(it, int64(-7354673213123121983), i64)

		err = Unmarshal(nullval, &i64)
		require.NoError(t, err)
		require.Equal(it, int64(0), i64)

		err = Unmarshal([]byte(`1354677198`), &i)
		require.NoError(t, err)
		require.Equal(it, int(1354677198), i)

		err = Unmarshal([]byte(`-1354677198`), &i)
		require.NoError(t, err)
		require.Equal(it, int(-1354677198), i)

		err = Unmarshal(nullval, &i)
		require.NoError(t, err)
		require.Equal(it, int(0), i)
	})

	t.Run("Unmarshal Unsigned Integer", func(it *testing.T) {
		var (
			ui8  uint8
			ui16 uint16
			ui32 uint32
			ui64 uint64
			ui   uint
		)

		err = Unmarshal([]byte(`10`), &ui8)
		require.NoError(t, err)
		require.Equal(it, uint8(10), ui8)

		err = Unmarshal([]byte(`-10`), &ui8)
		require.Error(t, err)

		err = Unmarshal(nullval, &ui8)
		require.NoError(t, err)
		require.Equal(it, uint8(0), ui8)

		err = Unmarshal([]byte(`128`), &ui16)
		require.NoError(t, err)
		require.Equal(it, uint16(128), ui16)

		err = Unmarshal([]byte(`-128`), &ui16)
		require.Error(t, err)

		err = Unmarshal(nullval, &ui16)
		require.NoError(t, err)
		require.Equal(it, uint16(0), ui16)

		err = Unmarshal([]byte(`1354677198`), &ui32)
		require.NoError(t, err)
		require.Equal(it, uint32(1354677198), ui32)

		err = Unmarshal([]byte(`-1354677198`), &ui32)
		require.Error(t, err)

		err = Unmarshal(nullval, &ui32)
		require.NoError(t, err)
		require.Equal(it, uint32(0), ui32)

		err = Unmarshal([]byte(`7354673213123121983`), &ui64)
		require.NoError(t, err)
		require.Equal(it, uint64(7354673213123121983), ui64)

		err = Unmarshal([]byte(`-7354673213123121983`), &ui64)
		require.Error(t, err)

		err = Unmarshal(nullval, &ui64)
		require.NoError(t, err)
		require.Equal(it, uint64(0), ui64)

		err = Unmarshal([]byte(`1354677198`), &ui)
		require.NoError(t, err)
		require.Equal(it, uint(1354677198), ui)

		err = Unmarshal([]byte(`-1354677198`), &ui)
		require.Error(t, err)

		err = Unmarshal(nullval, &ui)
		require.NoError(t, err)
		require.Equal(it, uint(0), ui)
	})

	t.Run("Unmarshal Float", func(it *testing.T) {
		var (
			f32 float32
			f64 float64
		)

		err = Unmarshal([]byte(`10`), &f32)
		require.NoError(t, err)
		require.Equal(it, float32(10), f32)

		err = Unmarshal([]byte(`10.32`), &f32)
		require.NoError(t, err)
		require.Equal(it, float32(10.32), f32)

		err = Unmarshal([]byte(`-882.3261239`), &f32)
		require.NoError(t, err)
		require.Equal(it, float32(-882.3261239), f32)

		err = Unmarshal([]byte(`-128.32128392`), &f64)
		require.NoError(t, err)
		require.Equal(it, float64(-128.32128392), f64)

		err = Unmarshal([]byte(`10.32128392`), &f64)
		require.NoError(t, err)
		require.Equal(it, float64(10.32128392), f64)
	})

	t.Run("Unmarshal Byte", func(it *testing.T) {
		b = []byte(`"` + string(byteval) + `"`)
		var bytea []byte
		err = Unmarshal(b, &bytea)
		require.NoError(t, err)
		require.Equal(t, pk, bytea)

		bytea = []byte(nil)
		err = Unmarshal(nullval, &bytea)
		require.NoError(t, err)
		require.Equal(t, []byte(nil), bytea)
	})

	t.Run("Unmarshal Time", func(it *testing.T) {
		var dt time.Time
		date := `2018-01-02T15:04:33Z`
		b = []byte(`"` + date + `"`)

		err = Unmarshal(b, &dt)
		require.NoError(t, err)
		require.Equal(t, date, dt.UTC().Format(time.RFC3339))

		err = Unmarshal(nullval, &dt)
		require.NoError(t, err)
		require.Equal(t, `0001-01-01T00:00:00Z`, dt.UTC().Format(time.RFC3339))
	})

	t.Run("Unmarshal Array", func(it *testing.T) {
		var (
			nullArr   []string
			initArr   []int
			strArr    []string
			intArr    []int
			twoDArr   [][]int
			threeDArr [][][]string
		)

		nullArr = []string{"xyz"}
		err = Unmarshal(nullval, &nullArr)
		require.NoError(t, err)
		require.Equal(t, []string(nil), nullArr)

		err = Unmarshal([]byte(`null`), &nullArr)
		require.NoError(t, err)
		require.Equal(t, []string(nil), nullArr)

		err = Unmarshal([]byte("[]"), &initArr)
		require.NoError(t, err)
		require.NotNil(t, initArr)
		require.Equal(t, make([]int, 0), initArr)

		err = Unmarshal([]byte(`["a", "b", "c"]`), &strArr)
		require.NoError(t, err)
		require.ElementsMatch(t, []string{"a", "b", "c"}, strArr)

		err = Unmarshal([]byte(`[2, 8, 32, 64, 128]`), &intArr)
		require.NoError(t, err)
		require.ElementsMatch(t, []int{2, 8, 32, 64, 128}, intArr)

		err = Unmarshal([]byte(`[
			[2, 8, 32, 64, 128],
			[1, 3, 5, 7],
			[0, 100, 1000, 10000, 100000]
		]`), &twoDArr)
		require.ElementsMatch(t, [][]int{
			{2, 8, 32, 64, 128},
			{1, 3, 5, 7},
			{0, 100, 1000, 10000, 100000},
		}, twoDArr)
		require.NoError(t, err)

		err = Unmarshal([]byte(`[
			[
				["a", "b", "c", "d", "e"],
				["甲", "乙", "丙", "丁"],
				["😄", "😅", "😆", "😉", "😊"]
			],
			[
				["a", "b", "c", "d", "e"],
				["f", "g", "h", "i", "j"]
			],
			[
				["Java", "JavaScript", "TypeScript"],
				["Rust", "GoLang"]
			]
		]`), &threeDArr)
		require.NoError(t, err)
		require.ElementsMatch(t, [][][]string{{
			[]string{"a", "b", "c", "d", "e"},
			[]string{"甲", "乙", "丙", "丁"},
			[]string{"😄", "😅", "😆", "😉", "😊"},
		}, {
			[]string{"a", "b", "c", "d", "e"},
			[]string{"f", "g", "h", "i", "j"},
		}, {
			[]string{"Java", "JavaScript", "TypeScript"},
			[]string{"Rust", "GoLang"},
		}}, threeDArr)

		b := []byte(`[{"Name":   "John", 
		"Email": "john@hotmail.com"   ,
		"Age" : 20,
		},  {"Name":   " 😃 😄 😅 😆 😉 😊 😋 😎 😍", 
		"Email": "test@hotmail.com"   ,
		"Age" : 88,
		"UUID": "daa68da0-8890-11ea-bc55-0242ac130003"
		}]`)

		users := []User{}
		err = Unmarshal(b, &users)
		require.NoError(t, err)
		id := uuid.MustParse("daa68da0-8890-11ea-bc55-0242ac130003")
		require.ElementsMatch(t, []User{
			{Name: "John", Email: "john@hotmail.com", Age: 20},
			{Name: " 😃 😄 😅 😆 😉 😊 😋 😎 😍", Email: "test@hotmail.com", Age: 88, UUID: id},
		}, users)

		type recursiveNestedStruct struct {
			Unicode   string
			String    string
			Boolean   bool
			Integer   int
			Level     uint
			Recursive *recursiveNestedStruct
			Float     float64
			Time      time.Time
			Nested    *struct {
				Name     string
				No       int
				StrSlice []string
			}
		}

		b = []byte(`
		[{
			"Unicode" : "中国★",
			"String": "hello world",
				"Level" : 1,
			"Time": "2017-02-05T03:46:15Z",
			"Recursive": {},
			"Nested": {
				"No" : 2,
				"StrSlice": ["a", "b", "c"],
			}
		}]
		`)

		structs := []recursiveNestedStruct{}
		err = Unmarshal(b, &structs)
		require.NoError(t, err)

		f := structs[0]
		require.NotNil(t, f.Nested)
		require.Equal(t, "中国★", f.Unicode)
		require.Equal(t, uint(1), f.Level)
		require.Equal(t, "hello world", f.String)
		require.Equal(t, 2, f.Nested.No)
		require.ElementsMatch(t, []string{"a", "b", "c"}, f.Nested.StrSlice)
	})

	t.Run("Unmarshal Map", func(it *testing.T) {
		data := make(map[customKey]string)

		err = Unmarshal([]byte(`{"test":"hello" ,  "test2": "world"}`), &data)
		require.NoError(t, err)

		k1 := customKey{value: "test"}
		k2 := customKey{value: "test2"}

		require.Contains(t, data, k1)
		require.Contains(t, data, k2)

		require.Equal(t, data[k1], "hello")
		require.Equal(t, data[k2], "world")
	})

	t.Run("Unmarshal Struct", func(it *testing.T) {
		{
			// unmarshal with empty object {}
			b = []byte(`   {   } `)
			var a struct{}
			err = Unmarshal(b, &a)
			require.NoError(it, err)
			require.Equal(t, struct{}{}, a)
		}

		{
			b = []byte(`
		{
			"Str" :"hello world!!" ,
			"UUID":     "4c03d1de-645b-40d2-9ed5-12bb537a602e",
			"SymbolStr"   : "x1#$%^\t!\n\t\t@#$%^&*())))?\\<>.,/:\":;'{}[]-=+_~",
			"EscapeStr"     :    "<html><div>hello world!</div></html>",
			"StrSlice" : ["a", "b", "c", "d"],
			"Users" : [
				{"Name":"SianLoong",   "Age": 18}   ,
			 { "Name":"Junkai"}],
			"Nested": {
				"NilKey": null,
				"MultiNilPtrKey": null,
				"MultiPtr": "testing \"multiple\" pointer",
				"Security"   : {
					"PrivateKey": "LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlDWFFJQkFBS0JnUUMrYUlTemtOdXFiZmdxWW9IYW1iS0dyaEF6UnV0dWYydWFzOUJIeXllUFJUdUk5ZVdwCnJHY3lRZlhPVlh2OGJBZVMxK2tIS0MvK1ZDTk9EbGZBTFlQWVVZa053eHVvRnFMbU1SR3E1MzMwSEVLSUpySDcKSUU5aUs0QUVZL3h5WjBTUEp5ZkNnQ2ZaeGtJTmpacWFoSS8rVWxrL1BmdWwyaEQ0ZTNUZVpGTm5HUUlEQVFBQgpBb0dBSHpOYlExMWlVV3dSdFVoTkJQZ1lnKzh6NG1NbG93YW9LRUo4eDdibmRaZWZxTks2WG5KTXVyU0tSZFJHCks5ZTc2ZmtOUzBudmkxcFlLcXM0LzltMWQ4Mk9XdmtDeXZvR3pmRXdyNGJ6bVBBZjdkczVkWElhb29wbWV4WWwKbFpsSmtuMDhWNFJmOWc4RFEyNFRsb3BpZ3RrSzY5UktRSzFHaHVyV1A4UjVxeTBDUVFEZ3dxcGVKZFF5RUdaYgpQUElJN2ZsSUVDRjlQNnNVc1ovNW40cEhNNmg2N0dOdGU1cEx4bDkzOXYybVhaN09aSUZHQU1rUmNDL1ZIK3c4Cm5oaytaNE9yQWtFQTJOK01oOWRpN1YreldaNUNIWXBWTHM5Qi9xOVl3YjFCNjN0UnZUbG9QSnFqTHc1NDUzZUUKbEs0ZnJSaVhXbEhLaUpLYlBOTU1ZUVkyTVRrcEQ2dDhTd0pCQUlkU2JRVFdQZFlPcmJITkZlUnVjeUlDSkVlbQpwN2lENFUrSDBOZGhzTlNoc3BOZVVkM0JpQVZRZmhOR1ZyRHBMalFaa1BXZzJBdTNkcUpnaGM1ZXdKVUNRQUVFCkV4RnoxZGZNMGZkQ2dZYkg1aHhCQmtzZUlTbFBMS2JndmdKSDZaQVhIVnFVRThicHpXb3c0cDhaOVdPTDdJbjEKUGRyc0ZpdkNMckRPVnIzbkRMOENRUURKSENwSEVFNTc0ckpzblJNYk5rR0F5dmZheW9MeElhVUF5WXovaGxrMgpzQ0wzb3BsdDNYM0tjYzV1MkRISVFsZTdGM1M4Wmp4REZMSVRrbnJ4QS9UVgotLS0tLUVORCBSU0EgUFJJVkFURSBLRVktLS0tLQ=="
				},
				"Test" :"hello world!!" ,
				"Test2" :"hello world!!" ,
				"Testxx" :"hello world!!" ,
				"empty" :    {},
				"nested"  : {
					"deep0"  : 100,
					"deep1" : {
						"value" : 199303.00
					},
					"deep2": "YOLO"
				}
			},
			"BigDecimal": 100.111,
			"Integer": 6000,
			"Bool": true
		}`)

			cp := make([]byte, len(b))
			copy(cp, b)

			var (
				o           testStruct
				nilKey      *types.Key
				multiNilKey ******types.Key
			)

			o.Nested.Key = types.IDKey("XX", 100, nil)

			err = Unmarshal(cp, &o)
			require.NoError(t, err)

			require.Nil(t, o.Nested.Key)
			require.Equal(t, nilKey, o.Nested.Key)

			require.Nil(t, o.Nested.NilKey)
			require.Equal(t, nilKey, o.Nested.NilKey)

			require.Nil(t, o.Nested.MultiNilPtrKey)
			require.Equal(t, multiNilKey, o.Nested.MultiNilPtrKey)

			// after unmarshal, the input should be the same (input shouldn't modified)
			require.Equal(t, uuid.MustParse(`4c03d1de-645b-40d2-9ed5-12bb537a602e`), *****o.UUID)
			require.Equal(t, `testing "multiple" pointer`, ***o.Nested.MultiPtr)
			require.Equal(t, b, cp)
			require.Equal(t, `hello world!!`, o.Str)
			require.Equal(t, `x1#$%^	!
		@#$%^&*())))?\<>.,/:":;'{}[]-=+_~`, o.SymbolStr)
			require.Equal(t, pk, o.Nested.Security.PrivateKey)
			require.Equal(t, true, o.Bool)
			require.Equal(t, int(6000), o.Integer)
			require.Equal(t, float64(100.111), o.BigDecimal)
			require.ElementsMatch(t, []User{
				{Name: "SianLoong", Age: 18},
				{Name: "Junkai"},
			}, o.Users)
			require.ElementsMatch(t, []string{"a", "b", "c", "d"}, o.StrSlice)
		}

		{
			var i User
			i.Name = "testing"
			i.Email = "sianloong90@gmail.com"
			i.Age = 100
			err = Unmarshal(nullval, &i)
			require.NoError(t, err)
			require.Equal(t, User{}, i)
		}

		{
			u := new(User)
			u.Name = "testing"
			err = Unmarshal([]byte(`{"Name": "lol", "Email":"test@hotmail.com", "Age": 18}`), u)
			require.NoError(t, err)
			require.Equal(t, "lol", u.Name)
			require.Equal(t, "test@hotmail.com", u.Email)
			require.Equal(t, int(18), u.Age)
		}
	})

	t.Run("Unmarshal Pointer Struct", func(it *testing.T) {
		var ptr *ptrStruct
		err = Unmarshal([]byte(`null`), &ptr)
		require.NoError(t, err)

		var nilptr *ptrStruct
		require.Equal(t, nilptr, ptr)

		initPtr := new(ptrStruct)
		err = Unmarshal([]byte(`{
			"PtrStr": "testing {}!@#$%^&*(\\",
			"PtrBool": true,
			"PtrInt": -100,
			"PtrInt8": -88,
			"PtrInt16": 8814,
			"PtrInt64": -88818111321351212,
			"PtrUint": 718222455,
			"PtrUint8": 173,
			"PtrUint16": 8814,
			"PtrUint32": 2031273814,
			"PtrUint64": 88818111321351212,
			"PtrStrSlice" : ["a", "b", "c", "d"],
			"PtrJSONRaw": {   "k1" : "value"  , "k2"  : "  value 1312$%^&*"},
			"PtrJSONNumber" : 183939.778,
			"PtrUintSlice" : [7, 88, 199, 200],
			"PtrF32Slice" :  [10.88, 728.99, 722.00],
			"PtrCurrencyUnit" : "MYR",
			"PtrLanguageTag" : "en",
			"PtrTime": "2008-01-26T00:00:58Z"
		}`), initPtr)
		require.NoError(t, err)

		{
			str := `testing {}!@#$%^&*(\`
			require.Equal(t, &str, initPtr.PtrStr)
			flag := true
			require.Equal(t, &flag, initPtr.PtrBool)
		}

		{
			i := int(-100)
			require.Equal(t, &i, initPtr.PtrInt)
			i8 := int8(-88)
			require.Equal(t, &i8, initPtr.PtrInt8)
			i16 := int16(8814)
			require.Equal(t, &i16, initPtr.PtrInt16)
			i64 := int64(-88818111321351212)
			require.Equal(t, &i64, initPtr.PtrInt64)
		}

		{
			ui := uint(718222455)
			require.Equal(t, &ui, initPtr.PtrUint)
			ui8 := uint8(173)
			require.Equal(t, &ui8, initPtr.PtrUint8)
			ui16 := uint16(8814)
			require.Equal(t, &ui16, initPtr.PtrUint16)
			ui32 := uint32(2031273814)
			require.Equal(t, &ui32, initPtr.PtrUint32)
			ui64 := uint64(88818111321351212)
			require.Equal(t, &ui64, initPtr.PtrUint64)
		}

		{
			var nilByte *[]byte
			require.Equal(t, nilByte, initPtr.PtrByte)
		}

		{
			raw := json.RawMessage(`{"k1":"value","k2":"  value 1312$%^&*"}`)
			require.Equal(t, &raw, initPtr.PtrJSONRaw)
		}

		{
			num := json.Number("183939.778")
			require.Equal(t, &num, initPtr.PtrJSONNumber)
			f64, _ := initPtr.PtrJSONNumber.Float64()
			require.Equal(t, float64(183939.778), f64)
			i64, _ := initPtr.PtrJSONNumber.Int64()
			require.Equal(t, int64(0), i64)
		}

		{
			slice := []string{"a", "b", "c", "d"}
			require.Equal(t, &slice, initPtr.PtrStrSlice)
		}

		{
			require.Nil(it, initPtr.PtrUUID)
			require.Nil(it, initPtr.PtrIntSlice)
		}

		{
			slice := []uint{7, 88, 199, 200}
			require.Equal(t, &slice, initPtr.PtrUintSlice)
		}

		{
			slice := []float32{10.88, 728.99, 722.00}
			require.Equal(t, &slice, initPtr.PtrF32Slice)
		}

		{
			lang := language.English
			require.Equal(t, &lang, initPtr.PtrLanguageTag)
		}

		{
			unit, _ := currency.ParseISO("MYR")
			require.Equal(t, &unit, initPtr.PtrCurrencyUnit)
		}

		{
			dt, _ := time.Parse(time.RFC3339, "2008-01-26T00:00:58Z")
			require.Equal(t, &dt, initPtr.PtrTime)
		}
	})
}

func BenchmarkJSONUnmarshal(b *testing.B) {
	data := []byte(`
	{         
		"Test" :"hello world!!" ,
		"Test2"   : "x1#$%^&*xx",
		"Test4": {
			"Test" :"hello world!!" ,
			"Test2" :"hello world!!" ,
			"Testxx" :"hello world!!" , 
			"empty" :    {},
			"nested"  : {
				"deep0"  : 100,
				"deep1" : {
					"value" : 199303.00
				},
				"deep2": "YOLO"
			}
		},
		"Test0": 100.111,
		"Test99": 6000,
		"Bool": true
	}   		
	
	`)
	var (
		o   testStruct
		err error
	)

	for n := 0; n < b.N; n++ {
		err = json.Unmarshal(data, &o)
		require.NoError(b, err)
	}
}

func BenchmarkJSONBUnmarshal(b *testing.B) {
	// data := []byte(`
	// {
	// 	"Test" :"hello world!!" ,
	// 	"Test2"   : "x1#$%^&*xx",
	// 	"Test4": {
	// 		"Test" :"hello world!!" ,
	// 		"Test2" :"hello world!!" ,
	// 		"Testxx" :"hello world!!" ,
	// 		"empty" :    {},
	// 		"nested"  : {
	// 			"deep0"  : 100,
	// 			"deep1" : {
	// 				"value" : 199303.00
	// 			},
	// 			"deep2": "YOLO"
	// 		}
	// 	},
	// 	"Test0": 100.111,
	// 	"Test99": 6000,
	// 	"Bool": true
	// }

	// `)
	// var (
	// 	o   testStruct
	// 	err error
	// )
	// for n := 0; n < b.N; n++ {
	// 	err = Unmarshal(data, &o)
	// 	require.NoError(b, err)
	// }
}
