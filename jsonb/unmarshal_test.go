package jsonb

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

type testStruct struct {
	Test  string
	Test0 float64
	Test2 string
	Test4 struct {
		Test   string
		Test2  string
		Nested struct {
			Deep struct {
				Value float64 `sqlike:"value"`
			} `sqlike:"deep1"`
			YOLO string `sqlike:"deep2"`
		} `sqlike:"nested"`
	}
	Test99 int
	Bool   bool
}

func TestUnUnmarshal(t *testing.T) {
	// var (
	// 	b []byte
	// 	// str string
	// 	err error
	// )

	// 	b = []byte(`"hello world\ta<html>"`)
	// 	err = UnUnmarshal(b, &str)
	// 	require.NoError(t, err)
	// 	require.Equal(t, `hello world	a<html>`, str)

	// 	b = []byte(`"LS0tLS1CRUdJTiBQVUJMSUMgS0VZLS0tLS0KTUlHZk1BMEdDU3FHU0liM0RRRUJBUVVBQTRHTkFEQ0JpUUtCZ1FDa2xRaW80VGVJWm82M1MwRnZOb25ZMi9uQQpaVXZybkRSUEl6RUtLNEE3SHU0VWp4TmhlYnh1RUEvUHFTSmd4T0lIVlBuQVNyU3dqK0lsUG9rY2RyUjZFa3luCjBjdmpqd2pHUnlBR2F3VmhmN1RXSGpreFRLNnBJSXFSaUJLNGgrRS9mUHdwdkpUaWVGQ1NtSVdvdlI4V3o2SnkKZUNucG1OclR6RzZaSmxKY3ZRSURBUUFCCi0tLS0tRU5EIFBVQkxJQyBLRVktLS0t"`)
	// 	var key []byte
	// 	err = UnUnmarshal(b, &key)
	// 	require.NoError(t, err)
	// 	require.Equal(t, []byte(`-----BEGIN PUBLIC KEY-----
	// MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCklQio4TeIZo63S0FvNonY2/nA
	// ZUvrnDRPIzEKK4A7Hu4UjxNhebxuEA/PqSJgxOIHVPnASrSwj+IlPokcdrR6Ekyn
	// 0cvjjwjGRyAGawVhf7TWHjkxTK6pIIqRiBK4h+E/fPwpvJTieFCSmIWovR8Wz6Jy
	// eCnpmNrTzG6ZJlJcvQIDAQAB
	// -----END PUBLIC KEY----`), key)

	// 	b = []byte(`100`)
	// 	var i16 int16
	// 	err = UnUnmarshal(b, &i16)
	// 	require.NoError(t, err)
	// 	require.Equal(t, int16(100), i16)

	// 	b = []byte(`200`)
	// 	var u16 uint16
	// 	err = UnUnmarshal(b, &u16)
	// 	require.NoError(t, err)
	// 	require.Equal(t, uint16(200), u16)

	// b = []byte(`{"message":"hello world!"}`)
	// var raw json.RawMessage
	// err = UnUnmarshal(b, &raw)
	// require.NoError(t, err)
	// 	require.Equal(t, json.RawMessage(b), raw)
	// require.Equal(t, json.RawMessage(b), raw)

	// b = []byte(`   	"abc   2131293"			`)
	// var str string
	// err = UnUnmarshal(b, &str)
	// log.Println("Final :", str)
	// require.NoError(t, err)

	// b = []byte(`{"message":"hello world!"}`)
	// var raw json.RawMessage
	// err = UnUnmarshal(b, &raw)
	// require.NoError(t, err)
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
