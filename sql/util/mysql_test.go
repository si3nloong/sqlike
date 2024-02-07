package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMySQLUtil(t *testing.T) {
	utl := MySQLUtil{}

	t.Run("TableName", func(t *testing.T) {
		require.Equal(t, "`db`.`table`", utl.TableName("db", "table"))
		require.Equal(t, "`db`.`A`", utl.TableName("db", "A"))
	})

	t.Run("Quote", func(t *testing.T) {
		require.Equal(t, "`abc`", utl.Quote("abc"))
	})

	t.Run("Var", func(t *testing.T) {
		require.Equal(t, "?", utl.Var(1))
		require.Equal(t, "?", utl.Var(10))
	})

	t.Run("Wrap", func(t *testing.T) {
		require.Equal(t, `'value'`, utl.Wrap("value"))
		require.Equal(t, `'10'`, utl.Wrap("10"))
	})

	t.Run("Regexp", func(t *testing.T) {
		require.True(t, sqlFuncRegexp.MatchString(`CURRENT_TIMESTAMP(6)`))
		require.False(t, sqlFuncRegexp.MatchString(`'CURRENT_TIMESTAMP(6)'`))
		require.False(t, sqlFuncRegexp.MatchString(`'AAA'`))
		require.False(t, sqlFuncRegexp.MatchString(`''`))
		require.False(t, sqlFuncRegexp.MatchString(`'0'`))
		require.False(t, sqlFuncRegexp.MatchString(`123`))
	})

	t.Run("WrapOnlyValue", func(t *testing.T) {
		require.Equal(t, `'value'`, utl.WrapOnlyValue("value"))
		require.Equal(t, `'10'`, utl.WrapOnlyValue("10"))
		require.Equal(t, `CURRENT_TIMESTAMP(10)`, utl.WrapOnlyValue("CURRENT_TIMESTAMP(10)"))
	})
}
