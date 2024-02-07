package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBytesToString(t *testing.T) {
	str := `hello world`
	require.Equal(t, str, BytesToString([]byte(str)))
}

func TestStringToBytes(t *testing.T) {
	bytes := []byte(`hello world`)
	require.Equal(t, bytes, StringToBytes(string(bytes)))
}

func TestString(t *testing.T) {
	msg := "hello world"
	blr := AcquireString()
	defer ReleaseString(blr)
	blr.WriteString(msg)

	require.Equal(t, msg, blr.String())
	blr.Reset()
	require.Equal(t, "", blr.String())
}

func TestUnsafeString(t *testing.T) {
	msg := `sqlike@1.6.0`
	b := []byte(msg)
	require.Equal(t, msg, UnsafeString(b))
}

func TestEscapeString(t *testing.T) {
	require.Equal(t, "sqlike", EscapeString("sqlike", '`'))
	require.Equal(t, "``sqlike``", EscapeString("`sqlike`", '`'))
}
