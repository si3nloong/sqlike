package sqldump

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestByteToString(t *testing.T) {
	require.Equal(t, `NULL`, byteToString([]byte(nil)))
	require.Equal(t, `""`, byteToString([]byte{}))
	require.Equal(t, `"abc"`, byteToString([]byte(`abc`)))
}

func TestNumToString(t *testing.T) {
	require.Equal(t, `10223.00`, numToString([]byte(`10223.00`)))
	require.Equal(t, `101`, numToString([]byte(`101`)))
}

func TestTimestampToString(t *testing.T) {
	require.Equal(t, `"2022-01-31 00:26:35"`, tsToString([]byte(`2022-01-31 00:26:35`)))
}

func TestDateToString(t *testing.T) {
	require.Equal(t, `"2022-01-31"`, dateToString([]byte(`2022-01-31`)))
	require.Equal(t, `"2022-02-01"`, dateToString([]byte(`2022-02-01`)))
}

func TestJsonToString(t *testing.T) {
	require.Equal(t, `"null"`, jsonToString([]byte(nil)))
	require.Equal(t, `"{}"`, jsonToString([]byte(`{}`)))
}
