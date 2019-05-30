package jsonb

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnmarshal(t *testing.T) {
	var (
		b   []byte
		str string
		err error
	)

	b = []byte(`"hello world\ta<html>"`)
	err = Unmarshal(b, &str)
	require.NoError(t, err)
	require.Equal(t, `hello world	a<html>`, str)

	b = []byte(`"LS0tLS1CRUdJTiBQVUJMSUMgS0VZLS0tLS0KTUlHZk1BMEdDU3FHU0liM0RRRUJBUVVBQTRHTkFEQ0JpUUtCZ1FDa2xRaW80VGVJWm82M1MwRnZOb25ZMi9uQQpaVXZybkRSUEl6RUtLNEE3SHU0VWp4TmhlYnh1RUEvUHFTSmd4T0lIVlBuQVNyU3dqK0lsUG9rY2RyUjZFa3luCjBjdmpqd2pHUnlBR2F3VmhmN1RXSGpreFRLNnBJSXFSaUJLNGgrRS9mUHdwdkpUaWVGQ1NtSVdvdlI4V3o2SnkKZUNucG1OclR6RzZaSmxKY3ZRSURBUUFCCi0tLS0tRU5EIFBVQkxJQyBLRVktLS0t"`)
	var key []byte
	err = Unmarshal(b, &key)
	require.NoError(t, err)
	require.Equal(t, []byte(`-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCklQio4TeIZo63S0FvNonY2/nA
ZUvrnDRPIzEKK4A7Hu4UjxNhebxuEA/PqSJgxOIHVPnASrSwj+IlPokcdrR6Ekyn
0cvjjwjGRyAGawVhf7TWHjkxTK6pIIqRiBK4h+E/fPwpvJTieFCSmIWovR8Wz6Jy
eCnpmNrTzG6ZJlJcvQIDAQAB
-----END PUBLIC KEY----`), key)

	b = []byte(`100`)
	var i16 int16
	err = Unmarshal(b, &i16)
	require.NoError(t, err)
	require.Equal(t, int16(100), i16)

	b = []byte(`200`)
	var u16 uint16
	err = Unmarshal(b, &u16)
	require.NoError(t, err)
	require.Equal(t, uint16(200), u16)

	b = []byte(`{"message":"hello world!"}`)
	var raw json.RawMessage
	err = Unmarshal(b, &raw)
	require.NoError(t, err)
	require.Equal(t, json.RawMessage(b), raw)

	b = []byte(`	
	[   
		"a",
		null, "b", "c"]`)
	var strArr []string
	err = Unmarshal(b, &strArr)
	require.NoError(t, err)
	// require.Equal(t, json.RawMessage(b), raw)
}
