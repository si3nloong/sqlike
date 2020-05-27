package mysql

import (
	"testing"

	sqlstmt "github.com/si3nloong/sqlike/sql/stmt"
	"github.com/stretchr/testify/require"
)

func TestUseDatabase(t *testing.T) {
	ms := New()
	stmt := sqlstmt.AcquireStmt(ms)
	defer sqlstmt.ReleaseStmt(stmt)
	ms.UseDatabase(stmt, "db")
	require.Equal(t, "USE `db`;", stmt.String())
	require.ElementsMatch(t, []interface{}{}, stmt.Args())
}

func TestGetDatabases(t *testing.T) {
	ms := New()
	stmt := sqlstmt.AcquireStmt(ms)
	defer sqlstmt.ReleaseStmt(stmt)
	ms.GetDatabases(stmt)
	require.Equal(t, "SHOW DATABASES;", stmt.String())
	require.ElementsMatch(t, []interface{}{}, stmt.Args())
}
