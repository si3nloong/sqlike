package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMySQLUtil(t *testing.T) {
	utl := MySQLUtil{}

	require.Equal(t, "`db`.`table`", utl.TableName("db", "table"))
	require.Equal(t, "`db`.`A`", utl.TableName("db", "A"))

	require.Equal(t, "`abc`", utl.Quote("abc"))
	require.Equal(t, "?", utl.Var(1))
	require.Equal(t, "?", utl.Var(10))
	require.Equal(t, `'value'`, utl.Wrap("value"))
	require.Equal(t, `'10'`, utl.Wrap("10"))

	require.True(t, sqlFuncRegexp.MatchString(`CURRENT_TIMESTAMP(6)`))
	require.False(t, sqlFuncRegexp.MatchString(`'CURRENT_TIMESTAMP(6)'`))
	require.False(t, sqlFuncRegexp.MatchString(`'AAA'`))
	require.False(t, sqlFuncRegexp.MatchString(`''`))
	require.False(t, sqlFuncRegexp.MatchString(`'0'`))
	require.False(t, sqlFuncRegexp.MatchString(`123`))
}
