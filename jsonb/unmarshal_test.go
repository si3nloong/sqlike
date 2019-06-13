package jsonb

import (
	"encoding/json"
	"io/ioutil"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type User struct {
	Name  string
	Email string
	Age   int
}

type testStruct struct {
	Str        string
	BigDecimal float64
	SymbolStr  string
	EscapeStr  string
	StrSlice   []string
	Users      []User
	Nested     struct {
		Security struct {
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

func TestUnmarshal(t *testing.T) {
	var (
		strval       = `hello world !!!!!!`
		symbolstrval = `\n\t\t\t<html>\n\t\t\t\t<div>Hello World</div>\n\t\t\t</html>\n\t\t\ttesting with !@#$%^&*(_)()()_()_((*??D|}A||||\\\\))\n\t\t`
		byteval      = []byte(`LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlDWFFJQkFBS0JnUUMrYUlTemtOdXFiZmdxWW9IYW1iS0dyaEF6UnV0dWYydWFzOUJIeXllUFJUdUk5ZVdwCnJHY3lRZlhPVlh2OGJBZVMxK2tIS0MvK1ZDTk9EbGZBTFlQWVVZa053eHVvRnFMbU1SR3E1MzMwSEVLSUpySDcKSUU5aUs0QUVZL3h5WjBTUEp5ZkNnQ2ZaeGtJTmpacWFoSS8rVWxrL1BmdWwyaEQ0ZTNUZVpGTm5HUUlEQVFBQgpBb0dBSHpOYlExMWlVV3dSdFVoTkJQZ1lnKzh6NG1NbG93YW9LRUo4eDdibmRaZWZxTks2WG5KTXVyU0tSZFJHCks5ZTc2ZmtOUzBudmkxcFlLcXM0LzltMWQ4Mk9XdmtDeXZvR3pmRXdyNGJ6bVBBZjdkczVkWElhb29wbWV4WWwKbFpsSmtuMDhWNFJmOWc4RFEyNFRsb3BpZ3RrSzY5UktRSzFHaHVyV1A4UjVxeTBDUVFEZ3dxcGVKZFF5RUdaYgpQUElJN2ZsSUVDRjlQNnNVc1ovNW40cEhNNmg2N0dOdGU1cEx4bDkzOXYybVhaN09aSUZHQU1rUmNDL1ZIK3c4Cm5oaytaNE9yQWtFQTJOK01oOWRpN1YreldaNUNIWXBWTHM5Qi9xOVl3YjFCNjN0UnZUbG9QSnFqTHc1NDUzZUUKbEs0ZnJSaVhXbEhLaUpLYlBOTU1ZUVkyTVRrcEQ2dDhTd0pCQUlkU2JRVFdQZFlPcmJITkZlUnVjeUlDSkVlbQpwN2lENFUrSDBOZGhzTlNoc3BOZVVkM0JpQVZRZmhOR1ZyRHBMalFaa1BXZzJBdTNkcUpnaGM1ZXdKVUNRQUVFCkV4RnoxZGZNMGZkQ2dZYkg1aHhCQmtzZUlTbFBMS2JndmdKSDZaQVhIVnFVRThicHpXb3c0cDhaOVdPTDdJbjEKUGRyc0ZpdkNMckRPVnIzbkRMOENRUURKSENwSEVFNTc0ckpzblJNYk5rR0F5dmZheW9MeElhVUF5WXovaGxrMgpzQ0wzb3BsdDNYM0tjYzV1MkRISVFsZTdGM1M4Wmp4REZMSVRrbnJ4QS9UVgotLS0tLUVORCBSU0EgUFJJVkFURSBLRVktLS0tLQ==`)
		nullval      = []byte(`null`)
	)

	var (
		o   testStruct
		pk  []byte
		b   []byte
		err error
	)

	pk, err = ioutil.ReadFile("./../examples/pk.pem")
	require.NoError(t, err)

	t.Run("Unmarshal String", func(it *testing.T) {
		var str string
		Unmarshal([]byte(`"`+strval+`"`), &str)
		require.Equal(it, strval, str)

		output := `
			<html>
				<div>Hello World</div>
			</html>
			testing with !@#$%^&*(_)()()_()_((*??D|}A||||\\))
		`

		var symbolstr string
		Unmarshal([]byte(`"`+symbolstrval+`"`), &symbolstr)
		require.Equal(it, output, symbolstr)

		Unmarshal(nullval, &str)
		require.Equal(it, "", str)
	})

	t.Run("Unmarshal Boolean", func(it *testing.T) {
		var flag bool
		Unmarshal([]byte(`true`), &flag)
		require.Equal(it, true, flag)
		Unmarshal([]byte(`false`), &flag)
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

		Unmarshal([]byte(`10`), &i8)
		require.Equal(it, int8(10), i8)
		Unmarshal([]byte(`-10`), &i8)
		require.Equal(it, int8(-10), i8)
		Unmarshal(nullval, &i8)
		require.Equal(it, int8(0), i8)

		Unmarshal([]byte(`128`), &i16)
		require.Equal(it, int16(128), i16)
		Unmarshal([]byte(`-128`), &i16)
		require.Equal(it, int16(-128), i16)
		Unmarshal(nullval, &i16)
		require.Equal(it, int16(0), i16)

		Unmarshal([]byte(`1354677198`), &i32)
		require.Equal(it, int32(1354677198), i32)
		Unmarshal([]byte(`-1354677198`), &i32)
		require.Equal(it, int32(-1354677198), i32)
		Unmarshal(nullval, &i32)
		require.Equal(it, int32(0), i32)

		Unmarshal([]byte(`7354673213123121983`), &i64)
		require.Equal(it, int64(7354673213123121983), i64)
		Unmarshal([]byte(`-7354673213123121983`), &i64)
		require.Equal(it, int64(-7354673213123121983), i64)
		Unmarshal(nullval, &i64)
		require.Equal(it, int64(0), i64)

		Unmarshal([]byte(`1354677198`), &i)
		require.Equal(it, int(1354677198), i)
		Unmarshal([]byte(`-1354677198`), &i)
		require.Equal(it, int(-1354677198), i)
		Unmarshal(nullval, &i)
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

		Unmarshal([]byte(`10`), &ui8)
		require.Equal(it, uint8(10), ui8)
		err = Unmarshal([]byte(`-10`), &ui8)
		require.Error(t, err)
		Unmarshal(nullval, &ui8)
		require.Equal(it, uint8(0), ui8)

		Unmarshal([]byte(`128`), &ui16)
		require.Equal(it, uint16(128), ui16)
		err = Unmarshal([]byte(`-128`), &ui16)
		require.Error(t, err)
		Unmarshal(nullval, &ui16)
		require.Equal(it, uint16(0), ui16)

		Unmarshal([]byte(`1354677198`), &ui32)
		require.Equal(it, uint32(1354677198), ui32)
		err = Unmarshal([]byte(`-1354677198`), &ui32)
		require.Error(t, err)
		Unmarshal(nullval, &ui32)
		require.Equal(it, uint32(0), ui32)

		Unmarshal([]byte(`7354673213123121983`), &ui64)
		require.Equal(it, uint64(7354673213123121983), ui64)
		err = Unmarshal([]byte(`-7354673213123121983`), &ui64)
		require.Error(t, err)
		Unmarshal(nullval, &ui64)
		require.Equal(it, uint64(0), ui64)

		Unmarshal([]byte(`1354677198`), &ui)
		require.Equal(it, uint(1354677198), ui)
		err = Unmarshal([]byte(`-1354677198`), &ui)
		require.Error(t, err)
		Unmarshal(nullval, &ui)
		require.Equal(it, uint(0), ui)
	})

	t.Run("Unmarshal Float", func(it *testing.T) {
		var (
			f32 float32
			f64 float64
		)

		Unmarshal([]byte(`10`), &f32)
		require.Equal(it, float32(10), f32)

		Unmarshal([]byte(`10.32`), &f32)
		require.Equal(it, float32(10.32), f32)

		Unmarshal([]byte(`-882.3261239`), &f32)
		require.Equal(it, float32(-882.3261239), f32)

		Unmarshal([]byte(`-128.32128392`), &f64)
		require.Equal(it, float64(-128.32128392), f64)

		Unmarshal([]byte(`10.32128392`), &f64)
		require.Equal(it, float64(10.32128392), f64)
	})

	t.Run("Unmarshal Byte", func(it *testing.T) {
		b = append(b, '"')
		b = append(b, byteval...)
		b = append(b, '"')
		var bytea []byte
		Unmarshal(b, &bytea)
		require.Equal(t, pk, bytea)

		bytea = []byte(nil)
		Unmarshal(nullval, &bytea)
		require.Equal(t, []byte(nil), bytea)
	})

	t.Run("Unmarshal Time", func(it *testing.T) {
		var dt time.Time
		date := `2018-01-02T15:04:33Z`
		b = []byte(`"` + date + `"`)

		Unmarshal(b, &dt)
		require.Equal(t, date, dt.Format(time.RFC3339))

		Unmarshal(nullval, &dt)
		require.Equal(t, `0001-01-01T00:00:00Z`, dt.Format(time.RFC3339))
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
		Unmarshal(nullval, &nullArr)
		require.Equal(t, []string(nil), nullArr)

		Unmarshal([]byte("[]"), &initArr)
		require.NotNil(t, initArr)
		require.Equal(t, make([]int, 0), initArr)

		Unmarshal([]byte(`["a", "b", "c"]`), &strArr)
		require.ElementsMatch(t, []string{"a", "b", "c"}, strArr)

		Unmarshal([]byte(`[2, 8, 32, 64, 128]`), &intArr)
		require.ElementsMatch(t, []int{2, 8, 32, 64, 128}, intArr)

		Unmarshal([]byte(`[
			[2, 8, 32, 64, 128],
			[1, 3, 5, 7],
			[0, 100, 1000, 10000, 100000]
		]`), &twoDArr)
		require.ElementsMatch(t, [][]int{
			[]int{2, 8, 32, 64, 128},
			[]int{1, 3, 5, 7},
			[]int{0, 100, 1000, 10000, 100000},
		}, twoDArr)

		Unmarshal([]byte(`[
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
	})

	t.Run("Unmarshal Struct", func(it *testing.T) {
		b = []byte(`
		{
			"Str" :"hello world!!" ,
			"SymbolStr"   : "x1#$%^\t!\n\t\t@#$%^&*())))?\\<>.,/:\":;'{}[]-=+_~",
			"EscapeStr"     :    "<html><div>hello world!</div></html>",
			"StrSlice" : ["a", "b", "c", "d"],
			"Users" : [
				{"Name":"SianLoong",   "Age": 18}   ,
			 { "Name":"Junkai"}],
			"Nested": {
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

		cp := make([]byte, len(b), len(b))
		copy(cp, b)

		err = Unmarshal(cp, &o)
		require.NoError(t, err)
		// after unmarshal, the input should be the same (input shouldn't modified)
		require.Equal(t, b, cp)

		require.Equal(t, `hello world!!`, o.Str)
		require.Equal(t, `x1#$%^	!
		@#$%^&*())))?\<>.,/:":;'{}[]-=+_~`, o.SymbolStr)
		require.Equal(t, pk, o.Nested.Security.PrivateKey)
		require.Equal(t, true, o.Bool)

		require.Equal(t, int(6000), o.Integer)
		require.Equal(t, float64(100.111), o.BigDecimal)
		require.ElementsMatch(t, []User{
			User{Name: "SianLoong", Age: 18},
			User{Name: "Junkai"},
		}, o.Users)
		require.ElementsMatch(t, []string{"a", "b", "c", "d"}, o.StrSlice)

		var o User
		o.Name = "testing"
		o.Email = "sianloong90@gmail.com"
		o.Age = 100
		Unmarshal(nullval, &o)
		require.Equal(t, User{}, o)

		var u User
		u.Name = "testing"
		Unmarshal([]byte(`{"Name": "lol", "Email":"test@hotmail.com", "Age": 18}`), &u)
		require.Equal(t, "lol", u.Name)
		require.Equal(t, "test@hotmail.com", u.Email)
		require.Equal(t, int(18), u.Age)
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
		err = Unmarshal(data, &o)
		require.NoError(b, err)
	}
}
