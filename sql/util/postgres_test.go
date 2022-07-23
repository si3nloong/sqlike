package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPostgresUtil(t *testing.T) {
	utl := PostgresUtil{}

	require.Equal(t, "`db`.`table`", utl.TableName("db", "table"))
	require.Equal(t, "`db`.`A`", utl.TableName("db", "A"))

	require.Equal(t, `"abc"`, utl.Quote("abc"))
	require.Equal(t, "$1", utl.Var(1))
	require.Equal(t, "$10", utl.Var(10))
	require.Equal(t, `'value'`, utl.Wrap("value"))
	require.Equal(t, `'10'`, utl.Wrap("10"))
}
