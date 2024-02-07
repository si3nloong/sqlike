package sql

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBaseSQL(t *testing.T) {
	s := CommonSql{}
	require.Equal(t, `"db"."table"`, s.TableName("db", "table"))
	require.Equal(t, "`1234`", s.Quote("1234"))
	require.Equal(t, "$10", s.Var(10))
	require.Equal(t, `'12Ab''cd'`, s.Wrap("12Ab'cd"))
}
