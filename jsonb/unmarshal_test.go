package jsonb

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"testing"

	"github.com/stretchr/testify/require"
)

type User struct {
	Name string
	Age  int
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
		strval  = `hello world !!!!!!`
		byteval = []byte(`LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlDWFFJQkFBS0JnUUMrYUlTemtOdXFiZmdxWW9IYW1iS0dyaEF6UnV0dWYydWFzOUJIeXllUFJUdUk5ZVdwCnJHY3lRZlhPVlh2OGJBZVMxK2tIS0MvK1ZDTk9EbGZBTFlQWVVZa053eHVvRnFMbU1SR3E1MzMwSEVLSUpySDcKSUU5aUs0QUVZL3h5WjBTUEp5ZkNnQ2ZaeGtJTmpacWFoSS8rVWxrL1BmdWwyaEQ0ZTNUZVpGTm5HUUlEQVFBQgpBb0dBSHpOYlExMWlVV3dSdFVoTkJQZ1lnKzh6NG1NbG93YW9LRUo4eDdibmRaZWZxTks2WG5KTXVyU0tSZFJHCks5ZTc2ZmtOUzBudmkxcFlLcXM0LzltMWQ4Mk9XdmtDeXZvR3pmRXdyNGJ6bVBBZjdkczVkWElhb29wbWV4WWwKbFpsSmtuMDhWNFJmOWc4RFEyNFRsb3BpZ3RrSzY5UktRSzFHaHVyV1A4UjVxeTBDUVFEZ3dxcGVKZFF5RUdaYgpQUElJN2ZsSUVDRjlQNnNVc1ovNW40cEhNNmg2N0dOdGU1cEx4bDkzOXYybVhaN09aSUZHQU1rUmNDL1ZIK3c4Cm5oaytaNE9yQWtFQTJOK01oOWRpN1YreldaNUNIWXBWTHM5Qi9xOVl3YjFCNjN0UnZUbG9QSnFqTHc1NDUzZUUKbEs0ZnJSaVhXbEhLaUpLYlBOTU1ZUVkyTVRrcEQ2dDhTd0pCQUlkU2JRVFdQZFlPcmJITkZlUnVjeUlDSkVlbQpwN2lENFUrSDBOZGhzTlNoc3BOZVVkM0JpQVZRZmhOR1ZyRHBMalFaa1BXZzJBdTNkcUpnaGM1ZXdKVUNRQUVFCkV4RnoxZGZNMGZkQ2dZYkg1aHhCQmtzZUlTbFBMS2JndmdKSDZaQVhIVnFVRThicHpXb3c0cDhaOVdPTDdJbjEKUGRyc0ZpdkNMckRPVnIzbkRMOENRUURKSENwSEVFNTc0ckpzblJNYk5rR0F5dmZheW9MeElhVUF5WXovaGxrMgpzQ0wzb3BsdDNYM0tjYzV1MkRISVFsZTdGM1M4Wmp4REZMSVRrbnJ4QS9UVgotLS0tLUVORCBSU0EgUFJJVkFURSBLRVktLS0tLQ==`)
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
		err = Unmarshal([]byte(`"`+strval+`"`), &str)
		require.NoError(it, err)
		require.Equal(it, strval, str)
	})

	t.Run("Unmarshal Boolean", func(it *testing.T) {
		var flag bool
		Unmarshal([]byte(`true`), &flag)
		require.Equal(it, true, flag)
		Unmarshal([]byte(`false`), &flag)
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

		Unmarshal([]byte(`128`), &i16)
		require.Equal(it, int16(128), i16)
		Unmarshal([]byte(`-128`), &i16)
		require.Equal(it, int16(-128), i16)

		Unmarshal([]byte(`1354677198`), &i32)
		require.Equal(it, int32(1354677198), i32)
		Unmarshal([]byte(`-1354677198`), &i32)
		require.Equal(it, int32(-1354677198), i32)

		Unmarshal([]byte(`7354673213123121983`), &i64)
		require.Equal(it, int64(7354673213123121983), i64)
		Unmarshal([]byte(`-7354673213123121983`), &i64)
		require.Equal(it, int64(-7354673213123121983), i64)

		Unmarshal([]byte(`1354677198`), &i)
		require.Equal(it, int(1354677198), i)
		Unmarshal([]byte(`-1354677198`), &i)
		require.Equal(it, int(-1354677198), i)
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

		Unmarshal([]byte(`128`), &ui16)
		require.Equal(it, uint16(128), ui16)
		err = Unmarshal([]byte(`-128`), &ui16)
		require.Error(t, err)

		Unmarshal([]byte(`1354677198`), &ui32)
		require.Equal(it, uint32(1354677198), ui32)
		err = Unmarshal([]byte(`-1354677198`), &ui32)
		require.Error(t, err)

		Unmarshal([]byte(`7354673213123121983`), &ui64)
		require.Equal(it, uint64(7354673213123121983), ui64)
		err = Unmarshal([]byte(`-7354673213123121983`), &ui64)
		require.Error(t, err)

		Unmarshal([]byte(`1354677198`), &ui)
		require.Equal(it, uint(1354677198), ui)
		err = Unmarshal([]byte(`-1354677198`), &ui)
		require.Error(t, err)
	})

	t.Run("Unmarshal Float", func(it *testing.T) {

	})

	t.Run("Unmarshal Byte", func(it *testing.T) {
		b = make([]byte, 0)
		Unmarshal(byteval, &b)
		log.Println(string(b))
		// require.Equal(t, pk, b)
	})

	t.Run("Unmarshal Array", func(it *testing.T) {
		var (
			strArr []string
			intArr []int
		)

		Unmarshal([]byte(`["a", "b", "c"]`), &strArr)
		require.ElementsMatch(t, []string{"a", "b", "c"}, strArr)

		Unmarshal([]byte(`[2, 8, 32, 64, 128]`), &intArr)
		require.ElementsMatch(t, []int{2, 8, 32, 64, 128}, intArr)
	})

	t.Run("Unmarshal Struct", func(it *testing.T) {
		b = []byte(`
		{
			"Str" :"hello world!!" ,
			"SymbolStr"   : "x1#$%^&*xx",
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
		require.Equal(t, `x1#$%^&*xx`, o.SymbolStr)
		require.Equal(t, pk, o.Nested.Security.PrivateKey)
		require.Equal(t, true, o.Bool)

		require.Equal(t, int(6000), o.Integer)
		require.Equal(t, float64(100.111), o.BigDecimal)
		require.ElementsMatch(t, []User{
			User{Name: "SianLoong", Age: 18},
			User{Name: "Junkai"},
		}, o.Users)
		require.ElementsMatch(t, []string{"a", "b", "c", "d"}, o.StrSlice)
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
