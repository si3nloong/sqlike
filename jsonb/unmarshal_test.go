package jsonb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshal(t *testing.T) {
	b := []byte(`"hello world\ta<html>"`)
	var str string
	Unmarshal(b, &str)
	assert.Equal(t, `hello world	a<html>`, str, "it should be match")

	b = []byte(`"LS0tLS1CRUdJTiBQVUJMSUMgS0VZLS0tLS0KTUlHZk1BMEdDU3FHU0liM0RRRUJBUVVBQTRHTkFEQ0JpUUtCZ1FDa2xRaW80VGVJWm82M1MwRnZOb25ZMi9uQQpaVXZybkRSUEl6RUtLNEE3SHU0VWp4TmhlYnh1RUEvUHFTSmd4T0lIVlBuQVNyU3dqK0lsUG9rY2RyUjZFa3luCjBjdmpqd2pHUnlBR2F3VmhmN1RXSGpreFRLNnBJSXFSaUJLNGgrRS9mUHdwdkpUaWVGQ1NtSVdvdlI4V3o2SnkKZUNucG1OclR6RzZaSmxKY3ZRSURBUUFCCi0tLS0tRU5EIFBVQkxJQyBLRVktLS0t"`)
	var key []byte
	Unmarshal(b, &key)
	assert.Equal(t, []byte(`-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCklQio4TeIZo63S0FvNonY2/nA
ZUvrnDRPIzEKK4A7Hu4UjxNhebxuEA/PqSJgxOIHVPnASrSwj+IlPokcdrR6Ekyn
0cvjjwjGRyAGawVhf7TWHjkxTK6pIIqRiBK4h+E/fPwpvJTieFCSmIWovR8Wz6Jy
eCnpmNrTzG6ZJlJcvQIDAQAB
-----END PUBLIC KEY----`), key, "it should be match")
}
