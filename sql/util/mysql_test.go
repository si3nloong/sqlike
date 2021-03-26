package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMySQLUtil(t *testing.T) {
	utl := MySQLUtil{}

	require.Equal(t, "`abc`", utl.Quote("abc"))
	require.Equal(t, "?", utl.Var(1))
	require.Equal(t, "?", utl.Var(10))
	require.Equal(t, `'value'`, utl.Wrap("value"))
}
